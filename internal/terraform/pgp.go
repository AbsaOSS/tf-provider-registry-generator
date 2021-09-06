package terraform

import (
	"bytes"
	"os"
	"strings"

	"github.com/k0da/tfreg-golang/internal/types"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
)

func findKey(keyring openpgp.EntityList, fpr []byte) *openpgp.Entity {
	for _, entity := range keyring {
		key := entity.PrimaryKey
		fpr := key.Fingerprint
		if key.Fingerprint == fpr {
			return entity
		}
	}
	return nil
}

func (p *Provider) ExtractPublicKey() (key *types.GPGPublicKey, err error) {
	// open global pubring.gpg
	pubring := "/Users/ab012ib/.gnupg/pubring.gpg"
	fpr := []byte(os.Getenv("GPG_FINGERPRINT"))
	keyringFileBuffer, _ := os.Open(pubring)
	defer keyringFileBuffer.Close()

	entityList, _ := openpgp.ReadKeyRing(keyringFileBuffer)
	pubkey := findKey(entityList, fpr)
	b := bytes.NewBuffer(nil)
	w, _ := armor.Encode(b, openpgp.PublicKeyType, nil)
	err = pubkey.Serialize(w)
	if err != nil {
		return nil, err
	}
	w.Close()
	ASCIIArmor := b.String()
	key = &types.GPGPublicKey{
		KeyID: pubkey.PrimaryKey.KeyIdString(),
		ASCIIArmor: strings.ReplaceAll(ASCIIArmor, "\n", "\\n"),
	}

	return
}
