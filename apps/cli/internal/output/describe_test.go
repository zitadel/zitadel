package output

import (
"bytes"
"strings"
"testing"
"time"

"google.golang.org/protobuf/types/known/timestamppb"

objectv2 "github.com/zitadel/zitadel/pkg/grpc/object/v2"
userpb "github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func TestDescribeMessage_User(t *testing.T) {
msg := &userpb.GetUserByIDResponse{
User: &userpb.User{
UserId:   "123",
Username: "hodor",
State:    userpb.UserState_USER_STATE_ACTIVE,
Type: &userpb.User_Human{
Human: &userpb.HumanUser{
Profile: &userpb.HumanProfile{
GivenName:  "hodor",
FamilyName: "hodor",
},
},
},
},
}

var buf bytes.Buffer
describeMessage(&buf, msg.ProtoReflect(), 0)
out := buf.String()

for _, want := range []string{"USER ID", "123", "USERNAME", "hodor", "USER_STATE_ACTIVE", "GIVEN NAME", "FAMILY NAME"} {
if !strings.Contains(out, want) {
t.Errorf("describe output missing %q\nOutput:\n%s", want, out)
}
}
}

func TestDescribeMessage_Timestamp(t *testing.T) {
fixedTime := time.Date(2024, 1, 5, 9, 12, 0, 0, time.UTC)
msg := &userpb.GetUserByIDResponse{
User: &userpb.User{
UserId: "abc",
Details: &objectv2.Details{
ChangeDate: timestamppb.New(fixedTime),
},
},
}

var buf bytes.Buffer
describeMessage(&buf, msg.ProtoReflect(), 0)
out := buf.String()

if !strings.Contains(out, "2024-01-05T09:12:00Z") {
t.Errorf("expected RFC3339 timestamp in output, got:\n%s", out)
}
}

func TestDescribeMessage_NilMessage(t *testing.T) {
msg := &userpb.GetUserByIDResponse{}
var buf bytes.Buffer
// Should not panic on empty/nil nested fields
describeMessage(&buf, msg.ProtoReflect(), 0)
}
