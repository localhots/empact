package db

import (
	"net/http"
	"strconv"
	"time"
)

type (
	StatRequest struct {
		Org  string
		Team string
		User string
		From int64
		To   int64
	}
	StatItem struct {
		Item  string `json:"item"`
		Value int    `json:"value"`
	}
	StatPoint struct {
		StatItem
		Timestamp uint64 `json:"ts"`
	}
)

const orgReposTopQuery = `
select
    c.repo as item,
    sum(c.commits) as value
from contribs c
join members m on
    c.author = m.user and
    c.owner = m.org
where
    m.id is not null and
    c.owner = ? and
    c.week >= ? and
    c.week <= ?
group by item
order by value desc`

const orgReposActivityQuery = `
select
    c.week as ts,
    c.repo as item,
    sum(c.commits) as value
from contribs c
join members m on
    c.author = m.user and
    c.owner = m.org
where
    m.id is not null and
    c.owner = ? and
    c.week >= ? and
    c.week <= ?
group by ts, item
order by ts, item`

const orgTeamsTopQuery = `
select
    t.name as item,
    sum(c.commits) value
from contribs c
join members m on
    c.author = m.user and
    c.owner = m.org
join teams t on
    m.team_id = t.id
where
    m.id is not null and
    c.owner = ? and
    c.week >= ? and
    c.week <= ?
group by item
order by value desc`

const orgTeamsActivityQuery = `
select
    c.week as ts,
    t.name as item,
    sum(c.commits) as value
from contribs c
join members m on
    c.author = m.user and
    c.owner = m.org
join teams t on
    m.team_id = t.id
where
    m.id is not null and
    c.owner = ? and
    c.week >= ? and
    c.week <= ?
group by ts, item
order by ts, item`

const orgUsersTopQuery = `
select
    c.author as item,
    sum(c.commits) value
from contribs c
join members m on
    c.author = m.user and
    c.owner = m.org
where
    m.id is not null and
    c.owner = ? and
    c.week >= ? and
    c.week <= ?
group by item
order by value desc`

func StatOrgReposTop(r *StatRequest) (res []StatItem) {
	defer measure("StatOrgReposTop", time.Now())
	mustSelect(&res, orgReposTopQuery, r.Org, r.From, r.To)
	return
}

func StatOrgReposActivity(r *StatRequest) (res []StatPoint) {
	defer measure("StatOrgReposActivity", time.Now())
	mustSelect(&res, orgReposActivityQuery, r.Org, r.From, r.To)
	return
}

func StatOrgTeamsTop(r *StatRequest) (res []StatItem) {
	defer measure("StatOrgTeamsTop", time.Now())
	mustSelect(&res, orgTeamsTopQuery, r.Org, r.From, r.To)
	return
}

func StatOrgTeamsActivity(r *StatRequest) (res []StatPoint) {
	defer measure("StatOrgTeamsActivity", time.Now())
	mustSelect(&res, orgTeamsActivityQuery, r.Org, r.From, r.To)
	return
}

func StatOrgUsersTop(r *StatRequest) (res []StatItem) {
	defer measure("StatOrgUsersTop", time.Now())
	mustSelect(&res, orgUsersTopQuery, r.Org, r.From, r.To)
	return
}

func ParseRequest(r *http.Request) *StatRequest {
	var err error
	var from, to int64
	if len(r.FormValue("from")) > 0 {
		if from, err = strconv.ParseInt(r.FormValue("from"), 10, 64); err != nil {
			panic(err)
		}
	}
	if len(r.FormValue("to")) > 0 {
		if to, err = strconv.ParseInt(r.FormValue("to"), 10, 64); err != nil {
			panic(err)
		}
	}
	return &StatRequest{
		Org:  r.FormValue("org"),
		Team: r.FormValue("team"),
		User: r.FormValue("user"),
		From: from,
		To:   to,
	}
}
