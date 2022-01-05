package storage

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/kralicky/opni-gateway/pkg/keyring"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/retry"
)

const (
	secretName = "opni-tenant-keyring"
)

type SecretStore struct {
	clientset *kubernetes.Clientset
	namespace string
}

var _ KeyringStore = (*SecretStore)(nil)

func NewInClusterSecretStore() *SecretStore {
	// check downwards api
	namespace, ok := os.LookupEnv("POD_NAMESPACE")
	if !ok {
		log.Fatal("POD_NAMESPACE environment variable not set")
	}
	rc, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err)
	}
	cs, err := kubernetes.NewForConfig(rc)
	if err != nil {
		log.Fatal(err)
	}
	return &SecretStore{
		clientset: cs,
		namespace: namespace,
	}
}

func (s *SecretStore) Put(ctx context.Context, kr keyring.Keyring) error {
	keyringData, err := kr.Marshal()
	if err != nil {
		return err
	}
	secret, err := s.clientset.CoreV1().
		Secrets(s.namespace).
		Get(ctx, secretName, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			newSecret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      secretName,
					Namespace: s.namespace,
				},
				Data: map[string][]byte{
					"keyring": keyringData,
				},
			}
			_, err := s.clientset.CoreV1().
				Secrets(s.namespace).
				Create(ctx, newSecret, metav1.CreateOptions{})
			return err
		} else {
			return err
		}
	} else {
		return retry.RetryOnConflict(retry.DefaultRetry, func() error {
			secret.Data["keyring"] = keyringData
			_, err = s.clientset.CoreV1().
				Secrets(s.namespace).
				Update(ctx, secret, metav1.UpdateOptions{})
			if err != nil {
				log.Println("error updating secret (will retry): ", err)
			}
			return err
		})
	}
}

func (s *SecretStore) Get(ctx context.Context) (keyring.Keyring, error) {
	secret, err := s.clientset.CoreV1().
		Secrets(s.namespace).
		Get(ctx, secretName, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if data, ok := secret.Data["keyring"]; ok {
		return keyring.Unmarshal(data)
	} else {
		return nil, errors.New("secret is missing keyring data")
	}
}
