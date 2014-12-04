// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package backups_test

import (
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"
	"gopkg.in/mgo.v2"

	"github.com/juju/juju/apiserver/backups"
	"github.com/juju/juju/apiserver/params"
)

func (s *backupsSuite) TestCreateOkay(c *gc.C) {
	s.PatchValue(backups.IsReady,
		func(*mgo.Session) (bool, error) { return true, nil },
	)
	s.setBackups(c, s.meta, "")
	var args params.BackupsCreateArgs
	result, err := s.api.Create(args)
	c.Assert(err, jc.ErrorIsNil)
	var expected params.BackupsMetadataResult
	expected.UpdateFromMetadata(s.meta)

	c.Check(result, gc.DeepEquals, expected)
}

func (s *backupsSuite) TestCreateNotes(c *gc.C) {
	s.PatchValue(backups.IsReady,
		func(*mgo.Session) (bool, error) { return true, nil },
	)
	s.meta.Notes = "this backup is important"
	s.setBackups(c, s.meta, "")
	args := params.BackupsCreateArgs{
		Notes: "this backup is important",
	}
	result, err := s.api.Create(args)
	c.Assert(err, jc.ErrorIsNil)
	var expected params.BackupsMetadataResult
	expected.UpdateFromMetadata(s.meta)
	expected.Notes = "this backup is important"

	c.Check(result, gc.DeepEquals, expected)
}

func (s *backupsSuite) TestCreateError(c *gc.C) {
	s.setBackups(c, nil, "failed!")
	s.PatchValue(backups.IsReady,
		func(*mgo.Session) (bool, error) { return true, nil },
	)
	var args params.BackupsCreateArgs
	_, err := s.api.Create(args)

	c.Check(err, gc.ErrorMatches, "failed!")
}
