package cmd

import (
	"fmt"
	"os"

	"github.com/gamer22026/cmsnr/pkg/deployment"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"sigs.k8s.io/yaml"
)

// deployCmd represents the cert command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "generates certificates for webhook server",
	Run:   generate,
}

func init() {
	serverCmd.AddCommand(deployCmd)
	deployCmd.Flags().String("version", Version, "version of the app")
	viper.BindPFlag("version", deployCmd.Flags().Lookup("version"))
}

func printKind(i interface{}) {
	o, err := yaml.Marshal(i)
	if err != nil {
		log.Errorf("error printing object: %s", err)
		os.Exit(2)
	}

	fmt.Printf("---\n%s\n", o)
}

func generate(cmd *cobra.Command, args []string) {
	mService := "cmsnr-mutating-webhook"
	vService := "cmsnr-validating-webhook"
	mSecret := fmt.Sprintf("mutating-%s", viper.GetString("secret"))
	vSecret := fmt.Sprintf("validating-%s", viper.GetString("secret"))
	port := viper.GetInt("port")
	registry := viper.GetString("registry")
	version := viper.GetString("version")

	mCert, mKey, err := deployment.GenerateCertificate(mService, namespace)
	if err != nil {
		log.Error(err)
		os.Exit(2)
	}

	vCert, vKey, err := deployment.GenerateCertificate(vService, namespace)
	if err != nil {
		log.Error(err)
		os.Exit(2)
	}

	mw := deployment.NewMutatingWebhookServer().NamespacedName(mService, namespace).MutatingWebhook(port, mCert).Rules()
	vw := deployment.NewValidatingWebhookServer().NamespacedName(vService, namespace).ValidatingWebhook(port, vCert).Rules()

	printKind(deployment.NewSA(namespace))
	printKind(deployment.NewClusterRole())
	printKind(deployment.NewClusterRolebinding(namespace))
	fmt.Println(deployment.NewCRD())
	printKind(deployment.NewDeployment(mService, namespace, registry, "mutating", mSecret, port, version))
	printKind(deployment.NewDeployment(vService, namespace, registry, "validating", vSecret, port, version))
	printKind(deployment.NewService(mService, namespace, port))
	printKind(deployment.NewService(vService, namespace, port))
	printKind(deployment.CertAsSecret(mCert, mKey, mSecret, namespace))
	printKind(deployment.CertAsSecret(vCert, vKey, vSecret, namespace))
	printKind(mw.Config)
	printKind(vw.Config)

}
