package encryption

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/k0da/tfreg-golang/internal/location"
	"github.com/k0da/tfreg-golang/internal/types"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
)

type IEncryption interface {
	ExtractPublicKey() (key *types.GPGPublicKey, err error)
}

type Gpg struct {
	location location.ILocation
}

func NewGpg(l location.ILocation) (g *Gpg, err error) {
	if l == nil {
		return g, fmt.Errorf("nil location")
	}
	g = &Gpg{
		location: l,
	}
	return
}

func (g *Gpg) ExtractPublicKey() (key *types.GPGPublicKey, err error) {
	keyringFile, err := os.Open(g.location.GPGPubring())
	if err != nil {
		return
	}
	defer keyringFile.Close()

	entityList, err := openpgp.ReadKeyRing(keyringFile)
	if err != nil {
		return
	}
	entity := findKey(entityList, g.location.GPGFingerprint())
	if entity == nil {
		return key, fmt.Errorf("nil entity for fingerprint %s", g.location.GPGFingerprint())
	}
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

func findKey(keyring openpgp.EntityList, fpr string) *openpgp.Entity {
	for _, entity := range keyring {
		key := entity.PrimaryKey
		s := fmt.Sprintf("%X", key.Fingerprint[:])
		if s == fpr {
			return entity
		}
	}
	return nil
}
