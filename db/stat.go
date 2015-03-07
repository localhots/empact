package db

import (
	"time"
)

type (
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

func StatOrgReposTop(org string, from, to int64) (res []StatItem) {
	defer measure("StatOrgReposTop", time.Now())
	mustSelect(&res, orgReposTopQuery, org, from, to)
	return
}

func StatOrgReposActivity(org string, from, to int64) (res []StatPoint) {
	defer measure("StatOrgReposActivity", time.Now())
	mustSelect(&res, orgReposActivityQuery, org, from, to)
	return
}

func StatOrgTeamsTop(org string, from, to int64) (res []StatItem) {
	defer measure("StatOrgTeamsTop", time.Now())
	mustSelect(&res, orgTeamsTopQuery, org, from, to)
	return
}

func StatOrgTeamsActivity(org string, from, to int64) (res []StatPoint) {
	defer measure("StatOrgTeamsActivity", time.Now())
	mustSelect(&res, orgTeamsActivityQuery, org, from, to)
	return
}

func StatOrgUsersTop(org string, from, to int64) (res []StatItem) {
	defer measure("StatOrgUsersTop", time.Now())
	mustSelect(&res, orgUsersTopQuery, org, from, to)
	return
}
