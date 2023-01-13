package reconciler

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"github.com/tektoncd/pipeline/test/diff"
)

func TestReportInfoFromRun(t *testing.T) {
	testCases := []struct {
		name string
		run  *v1beta1.CustomRun
		info *ReportInfo
		err  string
	}{
		{
			name: "valid run",
			run: &v1beta1.CustomRun{
				Spec: v1beta1.CustomRunSpec{
					Params: []v1beta1.Param{
						{
							Name:  repoKey,
							Value: *v1beta1.NewStructuredValues("some-org/some-repo"),
						}, {
							Name:  prNumberKey,
							Value: *v1beta1.NewStructuredValues("5"),
						}, {
							Name:  shaKey,
							Value: *v1beta1.NewStructuredValues("abcd1234"),
						}, {
							Name:  jobNameKey,
							Value: *v1beta1.NewStructuredValues("some-job"),
						}, {
							Name:  resultKey,
							Value: *v1beta1.NewStructuredValues("success"),
						}, {
							Name:  optionalKey,
							Value: *v1beta1.NewStructuredValues("false"),
						}, {
							Name:  logURLKey,
							Value: *v1beta1.NewStructuredValues("http://some/where"),
						},
					},
				},
			},
			info: &ReportInfo{
				Repo:       "some-org/some-repo",
				PRNumber:   5,
				SHA:        "abcd1234",
				JobName:    "some-job",
				Result:     "success",
				LogURL:     "http://some/where",
				IsOptional: false,
			},
		}, {
			name: "missing repo",
			run: &v1beta1.CustomRun{
				Spec: v1beta1.CustomRunSpec{
					Params: []v1beta1.Param{
						{
							Name:  prNumberKey,
							Value: *v1beta1.NewStructuredValues("5"),
						}, {
							Name:  shaKey,
							Value: *v1beta1.NewStructuredValues("abcd1234"),
						}, {
							Name:  jobNameKey,
							Value: *v1beta1.NewStructuredValues("some-job"),
						}, {
							Name:  resultKey,
							Value: *v1beta1.NewStructuredValues("success"),
						}, {
							Name:  optionalKey,
							Value: *v1beta1.NewStructuredValues("false"),
						}, {
							Name:  logURLKey,
							Value: *v1beta1.NewStructuredValues("http://some/where"),
						},
					},
				},
			},
			err: "missing field(s): repo",
		}, {
			name: "missing PR number",
			run: &v1beta1.CustomRun{
				Spec: v1beta1.CustomRunSpec{
					Params: []v1beta1.Param{
						{
							Name:  repoKey,
							Value: *v1beta1.NewStructuredValues("some-org/some-repo"),
						}, {
							Name:  shaKey,
							Value: *v1beta1.NewStructuredValues("abcd1234"),
						}, {
							Name:  jobNameKey,
							Value: *v1beta1.NewStructuredValues("some-job"),
						}, {
							Name:  resultKey,
							Value: *v1beta1.NewStructuredValues("success"),
						}, {
							Name:  optionalKey,
							Value: *v1beta1.NewStructuredValues("false"),
						}, {
							Name:  logURLKey,
							Value: *v1beta1.NewStructuredValues("http://some/where"),
						},
					},
				},
			},
			err: "missing field(s): prNumber",
		}, {
			name: "missing SHA",
			run: &v1beta1.CustomRun{
				Spec: v1beta1.CustomRunSpec{
					Params: []v1beta1.Param{
						{
							Name:  repoKey,
							Value: *v1beta1.NewStructuredValues("some-org/some-repo"),
						}, {
							Name:  prNumberKey,
							Value: *v1beta1.NewStructuredValues("5"),
						}, {
							Name:  jobNameKey,
							Value: *v1beta1.NewStructuredValues("some-job"),
						}, {
							Name:  resultKey,
							Value: *v1beta1.NewStructuredValues("success"),
						}, {
							Name:  optionalKey,
							Value: *v1beta1.NewStructuredValues("false"),
						}, {
							Name:  logURLKey,
							Value: *v1beta1.NewStructuredValues("http://some/where"),
						},
					},
				},
			},
			err: "missing field(s): sha",
		}, {
			name: "missing job name",
			run: &v1beta1.CustomRun{
				Spec: v1beta1.CustomRunSpec{
					Params: []v1beta1.Param{
						{
							Name:  repoKey,
							Value: *v1beta1.NewStructuredValues("some-org/some-repo"),
						}, {
							Name:  prNumberKey,
							Value: *v1beta1.NewStructuredValues("5"),
						}, {
							Name:  shaKey,
							Value: *v1beta1.NewStructuredValues("abcd1234"),
						}, {
							Name:  resultKey,
							Value: *v1beta1.NewStructuredValues("success"),
						}, {
							Name:  optionalKey,
							Value: *v1beta1.NewStructuredValues("false"),
						}, {
							Name:  logURLKey,
							Value: *v1beta1.NewStructuredValues("http://some/where"),
						},
					},
				},
			},
			err: "missing field(s): jobName",
		}, {
			name: "missing result",
			run: &v1beta1.CustomRun{
				Spec: v1beta1.CustomRunSpec{
					Params: []v1beta1.Param{
						{
							Name:  repoKey,
							Value: *v1beta1.NewStructuredValues("some-org/some-repo"),
						}, {
							Name:  prNumberKey,
							Value: *v1beta1.NewStructuredValues("5"),
						}, {
							Name:  shaKey,
							Value: *v1beta1.NewStructuredValues("abcd1234"),
						}, {
							Name:  jobNameKey,
							Value: *v1beta1.NewStructuredValues("some-job"),
						}, {
							Name:  optionalKey,
							Value: *v1beta1.NewStructuredValues("false"),
						}, {
							Name:  logURLKey,
							Value: *v1beta1.NewStructuredValues("http://some/where"),
						},
					},
				},
			},
			err: "missing field(s): result",
		}, {
			name: "non-string value",
			run: &v1beta1.CustomRun{
				Spec: v1beta1.CustomRunSpec{
					Params: []v1beta1.Param{{
						Name:  repoKey,
						Value: *v1beta1.NewStructuredValues("bob", "steve"),
					}, {
						Name:  prNumberKey,
						Value: *v1beta1.NewStructuredValues("5"),
					}, {
						Name:  shaKey,
						Value: *v1beta1.NewStructuredValues("abcd1234"),
					}, {
						Name:  jobNameKey,
						Value: *v1beta1.NewStructuredValues("some-job"),
					}, {
						Name:  resultKey,
						Value: *v1beta1.NewStructuredValues("success"),
					}, {
						Name:  optionalKey,
						Value: *v1beta1.NewStructuredValues("false"),
					}, {
						Name:  logURLKey,
						Value: *v1beta1.NewStructuredValues("http://some/where"),
					}},
				},
			},
			err: "invalid value: should be a string, is array: repo",
		}, {
			name: "non-int value",
			run: &v1beta1.CustomRun{
				Spec: v1beta1.CustomRunSpec{
					Params: []v1beta1.Param{{
						Name:  repoKey,
						Value: *v1beta1.NewStructuredValues("some-org/some-repo"),
					}, {
						Name:  prNumberKey,
						Value: *v1beta1.NewStructuredValues("five"),
					}, {
						Name:  shaKey,
						Value: *v1beta1.NewStructuredValues("abcd1234"),
					}, {
						Name:  jobNameKey,
						Value: *v1beta1.NewStructuredValues("some-job"),
					}, {
						Name:  resultKey,
						Value: *v1beta1.NewStructuredValues("success"),
					}, {
						Name:  optionalKey,
						Value: *v1beta1.NewStructuredValues("false"),
					}, {
						Name:  logURLKey,
						Value: *v1beta1.NewStructuredValues("http://some/where"),
					}},
				},
			},
			err: "invalid value: five should be a number: prNumber",
		}, {
			name: "non-bool value",
			run: &v1beta1.CustomRun{
				Spec: v1beta1.CustomRunSpec{
					Params: []v1beta1.Param{{
						Name:  repoKey,
						Value: *v1beta1.NewStructuredValues("some-org/some-repo"),
					}, {
						Name:  prNumberKey,
						Value: *v1beta1.NewStructuredValues("5"),
					}, {
						Name:  shaKey,
						Value: *v1beta1.NewStructuredValues("abcd1234"),
					}, {
						Name:  jobNameKey,
						Value: *v1beta1.NewStructuredValues("some-job"),
					}, {
						Name:  resultKey,
						Value: *v1beta1.NewStructuredValues("success"),
					}, {
						Name:  optionalKey,
						Value: *v1beta1.NewStructuredValues("banana"),
					}, {
						Name:  logURLKey,
						Value: *v1beta1.NewStructuredValues("http://some/where"),
					}},
				},
			},
			err: "invalid value: banana should be a bool: isOptional",
		}, {
			name: "invalid result value",
			run: &v1beta1.CustomRun{
				Spec: v1beta1.CustomRunSpec{
					Params: []v1beta1.Param{{
						Name:  repoKey,
						Value: *v1beta1.NewStructuredValues("some-org/some-repo"),
					}, {
						Name:  prNumberKey,
						Value: *v1beta1.NewStructuredValues("5"),
					}, {
						Name:  shaKey,
						Value: *v1beta1.NewStructuredValues("abcd1234"),
					}, {
						Name:  jobNameKey,
						Value: *v1beta1.NewStructuredValues("some-job"),
					}, {
						Name:  resultKey,
						Value: *v1beta1.NewStructuredValues("banana"),
					}, {
						Name:  optionalKey,
						Value: *v1beta1.NewStructuredValues("false"),
					}, {
						Name:  logURLKey,
						Value: *v1beta1.NewStructuredValues("http://some/where"),
					}},
				},
			},
			err: "invalid value: should be one of 'pending', 'success', or 'failure', but is 'banana': result",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			info, err := ReportInfoFromRun(tc.run)
			if err != nil {
				if tc.err == "" {
					t.Fatalf("expected no error, but got '%s'", err.Error())
				} else if tc.err != err.Error() {
					t.Fatalf("expected error '%s', but got '%s'", tc.err, err.Error())
				}
			} else {
				if tc.err != "" {
					t.Fatalf("expected error '%s', but got no error", tc.err)
				}

				if d := cmp.Diff(tc.info, info); d != "" {
					t.Errorf("result differs: %s", diff.PrintWantGot(d))
				}
			}
		})
	}
}
