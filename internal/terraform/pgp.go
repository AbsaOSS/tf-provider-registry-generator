package terraform

import (
	"bytes"
	"os"
	"strings"

	"github.com/k0da/tfreg-golang/internal/types"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
)

func findKey(keyring openpgp.EntityList, fpr string) *openpgp.Entity {
	for _, entity := range keyring {
		key := entity.PrimaryKey
		if string(key.Fingerprint[:]) == fpr {
			return entity
		}
	}
	return nil
}

func (p *Provider) ExtractPublicKey() (key *types.GPGPublicKey, err error) {
	keyringFile, err := os.Open(p.location.GPGPubring())
	if err != nil {
		return
	}
	defer keyringFile.Close()

	entityList, err := openpgp.ReadKeyRing(keyringFile)
	if err != nil {
		return
	}
	entity := findKey(entityList, p.location.GPGFingerprint())
	b := bytes.NewBuffer(nil)
	w, err := armor.Encode(b, openpgp.PublicKeyType, nil)
	if err != nil {
		return
	}
	defer w.Close()
	err = entity.Serialize(w)
	if err != nil {
		return nil, err
	}

	ASCIIArmor := b.String()
	key = &types.GPGPublicKey{
		KeyID:      entity.PrimaryKey.KeyIdString(),
		ASCIIArmor: strings.ReplaceAll(ASCIIArmor, "\n", "\\n"),
	}
	return
}
