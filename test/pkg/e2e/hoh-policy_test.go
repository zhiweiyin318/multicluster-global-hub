package tests

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog"
	policiesv1 "open-cluster-management.io/governance-policy-propagator/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/stolostron/multicluster-global-hub/test/pkg/utils"
)

const (
	INFORM_POLICY_YAML  = "../../resources/policy/inform-limitrange-policy.yaml"
	ENFORCE_POLICY_YAML = "../../resources/policy/enforce-limitrange-policy.yaml"

	POLICY_LABEL_KEY   = "global-policy"
	POLICY_LABEL_VALUE = "test"
	POLICY_NAME        = "policy-limitrange"
	POLICY_NAMESPACE   = "default"
)

var _ = Describe("Apply policy to the managed clusters", Ordered, Label("e2e-tests-policy"), func() {
	var token string
	var httpClient *http.Client
	var managedClusterName1 string
	var managedClusterName2 string
	var appClient client.Client

	BeforeAll(func() {
		By("Get token for the non-k8s-api")
		initToken, err := utils.FetchBearerToken(testOptions)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(len(initToken)).Should(BeNumerically(">", 0))
		token = initToken

		By("Config request of the api")
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		httpClient = &http.Client{Timeout: time.Second * 10, Transport: transport}

		By("Get managed cluster name")
		Eventually(func() error {
			managedClusters, err := getManagedCluster(httpClient, token)
			if err != nil {
				return err
			}
			managedClusterName1 = managedClusters[0].Name
			managedClusterName2 = managedClusters[1].Name
			return nil
		}, 3*time.Minute, 5*time.Second).ShouldNot(HaveOccurred())

		By("Get the appsubreport client")
		scheme := runtime.NewScheme()
		policiesv1.AddToScheme(scheme)
		appClient, err = clients.ControllerRuntimeClient(scheme)
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("add the label to a managedcluster for the policy", func() {
		patches := []patch{
			{
				Op:    "add", // or remove
				Path:  "/metadata/labels/" + POLICY_LABEL_KEY,
				Value: POLICY_LABEL_VALUE,
			},
		}
		Eventually(func() error {
			err := updateClusterLabel(httpClient, patches, token, managedClusterName1)
			if err != nil {
				return err
			}
			return nil
		}, 1*time.Minute, 1*time.Second).ShouldNot(HaveOccurred())

		By("Check the label is added")
		Eventually(func() error {
			managedCluster, err := getManagedClusterByName(httpClient, token, managedClusterName1)
			if err != nil {
				return err
			}
			if val, ok := managedCluster.Labels[POLICY_LABEL_KEY]; ok {
				if val == POLICY_LABEL_VALUE && managedCluster.Name == managedClusterName1 {
					return nil
				}
			}
			return fmt.Errorf("the label %s: %s is not exist", POLICY_LABEL_KEY, POLICY_LABEL_VALUE)
		}, 3*time.Minute, 5*time.Second).ShouldNot(HaveOccurred())
	})

	It("create a inform policy for the labeled cluster", func() {
		Eventually(func() error {
			_, err := clients.Kubectl(clients.HubClusterName(), "apply", "-f", INFORM_POLICY_YAML)
			if err != nil {
				return err
			}
			return nil
		}, 1*time.Minute, 1*time.Second).ShouldNot(HaveOccurred())

		Eventually(func() error {
			status, err := getPolicyStatus(appClient, POLICY_NAME, POLICY_NAMESPACE)
			if err != nil {
				return err
			}
			for _, policyInfo := range status.Status {
				if policyInfo.ClusterName == managedClusterName1 {
					if policyInfo.ComplianceState == policiesv1.NonCompliant {
						return nil
					}
				}
			}
			return fmt.Errorf("the policy have not applied to the managed cluster %s", managedClusterName1)
		}, 3*time.Minute, 5*time.Second).ShouldNot(HaveOccurred())
	})

	It("enforce the inform policy", func() {
		Eventually(func() error {
			_, err := clients.Kubectl(clients.HubClusterName(), "apply", "-f", ENFORCE_POLICY_YAML)
			if err != nil {
				return err
			}
			return nil
		}, 1*time.Minute, 1*time.Second).ShouldNot(HaveOccurred())

		Eventually(func() error {
			status, err := getPolicyStatus(appClient, POLICY_NAME, POLICY_NAMESPACE)
			if err != nil {
				return err
			}
			for _, policyInfo := range status.Status {
				if policyInfo.ClusterName == managedClusterName1 {
					if policyInfo.ComplianceState == policiesv1.Compliant {
						return nil
					}
				}
			}
			return fmt.Errorf("the policy have not applied to the managed cluster %s", managedClusterName1)
		}, 3*time.Minute, 5*time.Second).ShouldNot(HaveOccurred())
	})

	It("add policy to managedcluster2 by adding label", func() {
		patches := []patch{
			{
				Op:    "add", // or remove
				Path:  "/metadata/labels/" + POLICY_LABEL_KEY,
				Value: POLICY_LABEL_VALUE,
			},
		}

		Eventually(func() error {
			err := updateClusterLabel(httpClient, patches, token, managedClusterName2)
			if err != nil {
				return err
			}
			return nil
		}, 1*time.Minute, 1*time.Second).ShouldNot(HaveOccurred())

		By("Check the label is added")
		Eventually(func() error {
			managedCluster, err := getManagedClusterByName(httpClient, token, managedClusterName2)
			if err != nil {
				return err
			}
			if val, ok := managedCluster.Labels[POLICY_LABEL_KEY]; ok {
				if val == POLICY_LABEL_VALUE && managedCluster.Name == managedClusterName2 {
					return nil
				}
			}
			return fmt.Errorf("the label %s: %s is not exist", POLICY_LABEL_KEY, POLICY_LABEL_VALUE)
		}, 3*time.Minute, 5*time.Second).ShouldNot(HaveOccurred())

		By("Check the policy is added")
		Eventually(func() error {
			status, err := getPolicyStatus(appClient, POLICY_NAME, POLICY_NAMESPACE)
			if err != nil {
				return err
			}
			for _, policyInfo := range status.Status {
				if policyInfo.ClusterName == managedClusterName2 {
					if policyInfo.ComplianceState == policiesv1.Compliant {
						return nil
					}
				}
			}
			return fmt.Errorf("the policy have not applied to the managed cluster %s", managedClusterName2)
		}, 3*time.Minute, 5*time.Second).ShouldNot(HaveOccurred())
	})

	It("remove managedcluster1 policy by deleting label", func() {
		By("remove the label from the managedclusterName1")
		patches := []patch{
			{
				Op:    "remove",
				Path:  "/metadata/labels/" + POLICY_LABEL_KEY,
				Value: POLICY_LABEL_VALUE,
			},
		}
		Eventually(func() error {
			err := updateClusterLabel(httpClient, patches, token, managedClusterName1)
			if err != nil {
				return err
			}
			return nil
		}, 1*time.Minute, 1*time.Second).ShouldNot(HaveOccurred())

		By("Check the policy is removed from the managedclusterName1")
		Eventually(func() error {
			status, err := getPolicyStatus(appClient, POLICY_NAME, POLICY_NAMESPACE)
			if err != nil {
				return err
			}
			for _, policyInfo := range status.Status {
				if policyInfo.ClusterName == managedClusterName1 {
					return fmt.Errorf("the cluster %s policy(%s)should be removed",
						managedClusterName1, POLICY_NAME)
				}
			}
			return nil
		}, 3*time.Minute, 5*time.Second).ShouldNot(HaveOccurred())
	})

	AfterAll(func() {
		By("Delete the enforced policy")
		_, err := clients.Kubectl(clients.HubClusterName(), "delete", "-f", ENFORCE_POLICY_YAML)
		Expect(err).ShouldNot(HaveOccurred())

		By("Delete the LimitRange CR from managedcluster1 and managedcluster2")
		deleteInfo, err := clients.Kubectl(managedClusterName1, "delete", "LimitRange", "container-mem-limit-range")
		Expect(err).ShouldNot(HaveOccurred())
		klog.V(5).Info(managedClusterName1, ": ", deleteInfo)

		deleteInfo, err = clients.Kubectl(managedClusterName2, "delete", "LimitRange", "container-mem-limit-range")
		Expect(err).ShouldNot(HaveOccurred())
		klog.V(5).Info(managedClusterName2, ": ", deleteInfo)

		// TODO: clyang82 should check the policy is deleted in managed clusters of the regional hub cluster

		By("Delete the label from managedcluster2")
		patches := []patch{
			{
				Op:    "remove",
				Path:  "/metadata/labels/" + POLICY_LABEL_KEY,
				Value: POLICY_LABEL_VALUE,
			},
		}
		Eventually(func() error {
			err := updateClusterLabel(httpClient, patches, token, managedClusterName2)
			if err != nil {
				return err
			}
			return nil
		}, 1*time.Minute, 1*time.Second).ShouldNot(HaveOccurred())
	})
})

func getPolicyStatus(appClient client.Client, name, namespace string) (*policiesv1.PolicyStatus, error) {
	policy := &policiesv1.Policy{}
	err := appClient.Get(context.TODO(), types.NamespacedName{Namespace: namespace, Name: name}, policy)
	if err != nil {
		return nil, err
	}
	return &policy.Status, nil
}
