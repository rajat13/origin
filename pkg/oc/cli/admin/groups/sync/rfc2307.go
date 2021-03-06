package sync

import (
	"github.com/openshift/library-go/pkg/security/ldapclient"
	ldapquery "github.com/openshift/library-go/pkg/security/ldapquery"
	"github.com/openshift/origin/pkg/cmd/server/apis/config"
	"github.com/openshift/origin/pkg/oc/lib/groupsync"
	"github.com/openshift/origin/pkg/oc/lib/groupsync/interfaces"
	"github.com/openshift/origin/pkg/oc/lib/groupsync/rfc2307"
	"github.com/openshift/origin/pkg/oc/lib/groupsync/syncerror"
)

var _ SyncBuilder = &RFC2307Builder{}
var _ PruneBuilder = &RFC2307Builder{}

type RFC2307Builder struct {
	ClientConfig ldapclient.Config
	Config       *config.RFC2307Config

	rfc2307LDAPInterface *rfc2307.LDAPInterface

	ErrorHandler syncerror.Handler
}

func (b *RFC2307Builder) GetGroupLister() (interfaces.LDAPGroupLister, error) {
	return b.getRFC2307LDAPInterface()
}

func (b *RFC2307Builder) GetGroupNameMapper() (interfaces.LDAPGroupNameMapper, error) {
	ldapInterface, err := b.getRFC2307LDAPInterface()
	if err != nil {
		return nil, err
	}
	if b.Config.GroupNameAttributes != nil {
		return syncgroups.NewEntryAttributeGroupNameMapper(b.Config.GroupNameAttributes, ldapInterface), nil
	}

	return nil, nil
}

func (b *RFC2307Builder) GetUserNameMapper() (interfaces.LDAPUserNameMapper, error) {
	return syncgroups.NewUserNameMapper(b.Config.UserNameAttributes), nil
}

func (b *RFC2307Builder) GetGroupMemberExtractor() (interfaces.LDAPMemberExtractor, error) {
	return b.getRFC2307LDAPInterface()
}

func (b *RFC2307Builder) getRFC2307LDAPInterface() (*rfc2307.LDAPInterface, error) {
	if b.rfc2307LDAPInterface != nil {
		return b.rfc2307LDAPInterface, nil
	}

	groupQuery, err := ldapquery.NewLDAPQueryOnAttribute(ToLDAPQuery(b.Config.AllGroupsQuery), b.Config.GroupUIDAttribute)
	if err != nil {
		return nil, err
	}
	userQuery, err := ldapquery.NewLDAPQueryOnAttribute(ToLDAPQuery(b.Config.AllUsersQuery), b.Config.UserUIDAttribute)
	if err != nil {
		return nil, err
	}
	b.rfc2307LDAPInterface = rfc2307.NewLDAPInterface(b.ClientConfig,
		groupQuery, b.Config.GroupNameAttributes, b.Config.GroupMembershipAttributes,
		userQuery, b.Config.UserNameAttributes, b.ErrorHandler)
	return b.rfc2307LDAPInterface, nil
}

func (b *RFC2307Builder) GetGroupDetector() (interfaces.LDAPGroupDetector, error) {
	return b.getRFC2307LDAPInterface()
}
