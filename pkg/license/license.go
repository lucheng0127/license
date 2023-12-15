package license

import (
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/7Linternational/dmidecode"
	"github.com/lucheng0127/license/pkg/cipher"
)

const (
	LICENSEFILE string = "LICENSE"
)

type LicenseManager interface {
	Import(string) error
	GenerateLicense(string, int) (string, error)
	ValidateLicense(string) (bool, error)
	ParseLicense(string) (string, int, error)
	LifeTime() (string, error)
}

type LicenseMgr struct {
	SecKey     string
	LicenseDir string
	Cipher     cipher.Cipher
}

func NewLicenseManager(key, dir string) (LicenseManager, error) {
	mgr := &LicenseMgr{SecKey: key, LicenseDir: dir}

	c, err := cipher.NewAESCipher(mgr.SecKey)
	if err != nil {
		return nil, err
	}

	mgr.Cipher = c
	return mgr, nil
}

func (mgr *LicenseMgr) Import(licenseStr string) error {
	filename := path.Join(mgr.LicenseDir, LICENSEFILE)

	if _, err := os.Stat(mgr.LicenseDir); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(mgr.LicenseDir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	fo, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fo.Close()

	_, err = fo.Write([]byte(licenseStr))
	if err != nil {
		return err
	}

	return nil
}

func (mgr *LicenseMgr) getDmicode() (string, error) {
	dmi := dmidecode.New()
	if err := dmi.Run(); err != nil {
		return "", errors.New("permission denied Can't read memory from /dev/mem")
	}

	records, err := dmi.SearchByType(1)
	if err != nil {
		return "", err
	}

	return records[0]["UUID"], nil
}

func (mgr *LicenseMgr) GenerateLicense(dmiCode string, day int) (string, error) {
	currentSec := time.Now().Unix()
	expireSec := currentSec + int64(day*24*60*60)

	rawLic := fmt.Sprintf("%s_%d", dmiCode, expireSec)

	lic, err := mgr.Cipher.Encrypt([]byte(rawLic))
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(lic), nil
}

func (mgr *LicenseMgr) ParseLicense(licStr string) (string, int, error) {
	licBytes, err := hex.DecodeString(licStr)
	if err != nil {
		return "", -1, err
	}

	lic, err := mgr.Cipher.Decrypt(licBytes)
	if err != nil {
		return "", -1, err
	}

	licStrArray := strings.Split(string(lic), "_")
	if len(licStrArray) != 2 {
		return "", -1, errors.New("invalidate licence")
	}

	deadline, err := strconv.Atoi(licStrArray[1])
	if err != nil {
		return "", -1, errors.New("invalidate licence")
	}

	return licStrArray[0], deadline, nil
}

func (mgr *LicenseMgr) ValidateLicense(licStr string) (bool, error) {
	licDmiCode, deadline, err := mgr.ParseLicense(licStr)
	if err != nil {
		return false, err
	}

	dmiCode, err := mgr.getDmicode()
	if err != nil {
		return false, fmt.Errorf("validate license error: %s", err.Error())
	}

	if dmiCode != licDmiCode {
		return false, errors.New("invalidate license for this hyper")
	}

	if time.Now().Unix() > int64(deadline) {
		return false, errors.New("licence expired")
	}

	return true, nil
}

func (mgr *LicenseMgr) LifeTime() (string, error) {
	filename := path.Join(mgr.LicenseDir, LICENSEFILE)

	lic, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	_, life, err := mgr.ParseLicense(string(lic))
	if err != nil {
		return "", err
	}

	lifeTime := time.Unix(int64(life), 0)
	return lifeTime.Format("2006-01-02"), nil
}
