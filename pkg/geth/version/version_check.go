// source: https://github.com/ethereum/go-ethereum/blob/v1.10.11/cmd/geth/version_check.go

package version

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/ethereum/go-ethereum/log"
	"github.com/jedisct1/go-minisign"
	"golang.org/x/sync/errgroup"
)

var gethPubKeys []string = []string{
	//@holiman, minisign public key FB1D084D39BAEC24
	"RWQk7Lo5TQgd+wxBNZM+Zoy+7UhhMHaWKzqoes9tvSbFLJYZhNTbrIjx",
	//minisign public key 138B1CA303E51687
	"RWSHFuUDoxyLEzjszuWZI1xStS66QTyXFFZG18uDfO26CuCsbckX1e9J",
	//minisign public key FD9813B2D2098484
	"RWSEhAnSshOY/b+GmaiDkObbCWefsAoavjoLcPjBo1xn71yuOH5I+Lts",
}

type vulnJson struct {
	Name        string
	Uid         string
	Summary     string
	Description string
	Links       []string
	Introduced  string
	Fixed       string
	Published   string
	Severity    string
	Check       string
	CVE         string
}

const VersionCheckUrl = "https://geth.ethereum.org/docs/vulnerabilities/vulnerabilities.json"

// CheckVersion checks the version of Geth for vulnerabilities.
// The current url must pass the regex check of the vulnerability json file.
// The current version must be in following format: `Geth/v1.2.3-`
func CheckVersion(url, current string) (*vulnJson, error) {
	var (
		data []byte
		sig  []byte
		err  error
	)
	if data, err = fetch(url); err != nil {
		return nil, fmt.Errorf("could not retrieve data: %w", err)
	}
	if sig, err = fetch(fmt.Sprintf("%v.minisig", url)); err != nil {
		return nil, fmt.Errorf("could not retrieve signature: %w", err)
	}
	if err = verifySignature(gethPubKeys, data, sig); err != nil {
		return nil, err
	}
	var vulns []vulnJson
	if err = json.Unmarshal(data, &vulns); err != nil {
		return nil, err
	}
	// allOk := true
	for _, vuln := range vulns {
		r, err := regexp.Compile(vuln.Check)
		if err != nil {
			return nil, err
		}
		if r.MatchString(current) {
			// allOk = false
			// fmt.Printf("## Vulnerable to %v (%v)\n\n", vuln.Uid, vuln.Name)
			// fmt.Printf("Severity: %v\n", vuln.Severity)
			// fmt.Printf("Summary : %v\n", vuln.Summary)
			// fmt.Printf("Fixed in: %v\n", vuln.Fixed)
			// if len(vuln.CVE) > 0 {
			// 	fmt.Printf("CVE: %v\n", vuln.CVE)
			// }
			// if len(vuln.Links) > 0 {
			// 	fmt.Printf("References:\n")
			// 	for _, ref := range vuln.Links {
			// 		fmt.Printf("\t- %v\n", ref)
			// 	}
			// }
			// fmt.Println()
			return &vuln, nil
		}
	}
	// if allOk {
	// 	fmt.Println("No vulnerabilities found")
	// }
	return nil, nil
}

// fetch makes an HTTP request to the given url and returns the response body
func fetch(url string) ([]byte, error) {
	if filep := strings.TrimPrefix(url, "file://"); filep != url {
		return ioutil.ReadFile(filep)
	}
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// verifySignature checks that the sigData is a valid signature of the given
// data, for pubkey GethPubkey
func verifySignature(pubkeys []string, data, sigdata []byte) error {
	sig, err := minisign.DecodeSignature(string(sigdata))
	if err != nil {
		return err
	}
	// find the used key
	var key *minisign.PublicKey
	for _, pubkey := range pubkeys {
		pub, err := minisign.NewPublicKey(pubkey)
		if err != nil {
			// our pubkeys should be parseable
			return err
		}
		if pub.KeyId != sig.KeyId {
			continue
		}
		key = &pub
		break
	}
	if key == nil {
		log.Info("Signing key not trusted", "keyid", keyID(sig.KeyId), "error", err)
		return errors.New("signature could not be verified")
	}
	if ok, err := key.Verify(data, sig); !ok || err != nil {
		log.Info("Verification failed error", "keyid", keyID(key.KeyId), "error", err)
		return errors.New("signature could not be verified")
	}
	return nil
}

// keyID turns a binary minisign key ID into a hex string.
// Note: key IDs are printed in reverse byte order.
func keyID(id [8]byte) string {
	var rev [8]byte
	for i := range id {
		rev[len(rev)-1-i] = id[i]
	}
	return fmt.Sprintf("%X", rev)
}

// SemanticVersionToGethVersion converts a semantic version to a Geth version - primarily to match the regex check
// when identifying vulnerabile versions.
// Converts: `1.10.0` -> `Geth/v1.10.0-`
func SemanticVersionToGethVersion(sv string) string {
	return fmt.Sprintf("Geth/v%s-", sv)
}

// GetVulnerabilities returns a list of vulnerabilities for specified Geth releases (in parallel)
// TODO: testing
func GetVulnerabilities(semanticVersions ...string) ([]vulnJson, error) {
	var vulnerabilityList []vulnJson
	var g errgroup.Group

	for _, version := range semanticVersions {
		version := version
		g.Go(func() error {
			gethVersion := SemanticVersionToGethVersion(version)
			vuln, err := CheckVersion(VersionCheckUrl, gethVersion)
			if err == nil {
				vulnerabilityList = append(vulnerabilityList, *vuln)
			}
			return err
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}

	return vulnerabilityList, nil
}
