package reconciler

import (
	"context"
	"fmt"

	"github.com/jenkins-x/go-scm/scm"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"knative.dev/pkg/logging"
	kreconciler "knative.dev/pkg/reconciler"
)

// Reconciler is the core of the implementation of the PR commenter, adding, updating, or deleting comments as needed.
type Reconciler struct {
	SCMClient *scm.Client
	BotUser   string
}

// ReconcileKind implements Interface.ReconcileKind.
func (c *Reconciler) ReconcileKind(ctx context.Context, r *v1beta1.CustomRun) kreconciler.Event {
	logger := logging.FromContext(ctx)
	logger.Infof("Reconciling %s/%s", r.Namespace, r.Name)

	// Ignore completed waits.
	if r.IsDone() {
		logger.Info("Run is finished, done reconciling")
		return nil
	}

	if r.Spec.CustomRef == nil ||
		r.Spec.CustomRef.APIVersion != "custom.tekton.dev/v0" || r.Spec.CustomRef.Kind != "PRStatusUpdater" {
		// This is not a Run we should have been notified about; do nothing.
		return nil
	}
	if r.Spec.CustomRef.Name != "" {
		r.Status.MarkCustomRunFailed("UnexpectedName", "Found unexpected ref name: %s", r.Spec.CustomRef.Name)
		return fmt.Errorf("unexpected ref name: %s", r.Spec.CustomRef.Name)
	}

	spec, fieldErr := StatusInfoFromRun(r)
	if fieldErr != nil {
		r.Status.MarkCustomRunFailed("InvalidParams", "Invalid parameters: %s", fieldErr.Error())
		return fieldErr
	}

	gitRepoStatus := &scm.StatusInput{
		State:  scm.ToState(spec.State),
		Label:  spec.JobName,
		Desc:   spec.Description,
		Target: spec.TargetURL,
	}

	logger.Infof("creating status on repo %s for sha %s: %+v", spec.Repo, spec.SHA, gitRepoStatus)
	_, resp, err := c.SCMClient.Repositories.CreateStatus(ctx, spec.Repo, spec.SHA, gitRepoStatus)
	if err != nil {
		logger.Errorf("failure in SCM client: error: %v, headers: %+v", err, resp.Header)
		r.Status.MarkCustomRunFailed("SCMError", "Error interacting with SCM: %s", err.Error())
		return err
	}

	r.Status.MarkCustomRunSucceeded("Commented", "PR status successfully set")

	// Don't emit events on nop-reconciliations, it causes scale problems.
	return nil
}
