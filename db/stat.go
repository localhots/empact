package db

import (
	"fmt"
	"time"
)

type (
	StatItem struct {
		Item    string `json:"item"`
		Commits int    `json:"commits"`
		Delta   int    `json:"delta"`
	}
	StatPoint struct {
		StatItem
		Week uint64 `json:"week"`
	}
)

const orgTopQuery = `
select
    %s as item,
    sum(c.commits) as commits,
    sum(c.additions) - sum(c.deletions) as delta
from contribs c
join members m on
    c.author = m.user and
    c.owner = m.org
join teams t on
    m.team_id = t.id
where
    m.id is not null and
    c.owner = :org and
    c.week >= :from and
    c.week <= :to
group by item
order by %s desc`

const orgActivityQuery = `
select
    %s as item,
    sum(c.commits) as commits,
    sum(c.additions) - sum(c.deletions) as delta,
    c.week as week
from contribs c
join members m on
    c.author = m.user and
    c.owner = m.org
join teams t on
    m.team_id = t.id
where
    m.id is not null and
    c.owner = :org and
    c.week >= :from and
    c.week <= :to
group by item, week
order by week, %s desc`

const teamTopQuery = `
select
    %s as item,
    sum(c.commits) as commits,
    sum(c.additions) - sum(c.deletions) as delta
from contribs c
join members m on
    c.author = m.user and
    c.owner = m.org
join teams t on
    m.team_id = t.id and
    t.name = :team
where
    m.id is not null and
    c.owner = :org and
    c.week >= :from and
    c.week <= :to
group by item
order by commits desc`

const teamActivityQuery = `
select
    %s as item,
    sum(c.commits) as commits,
    sum(c.additions) - sum(c.deletions) as delta,
    c.week as week
from contribs c
join members m on
    c.author = m.user and
    c.owner = m.org
join teams t on
    m.team_id = t.id and
    t.name = :team
where
    m.id is not null and
    c.owner = :org and
    c.week >= :from and
    c.week <= :to
group by item, week
order by week, commits desc`

const userTopQuery = `
select
    c.repo as item,
    sum(c.commits) as commits,
    sum(c.additions) - sum(c.deletions) as delta
from contribs c
where
    c.owner = :org and
    c.week >= :from and
    c.week <= :to
group by item
order by commits desc`

const userActivityQuery = `
select
    c.week as week,
    c.repo as item,
    sum(c.commits) as commits,
    sum(c.additions) - sum(c.deletions) as delta
from contribs c
where
    c.owner = :org and
    c.week >= :from and
    c.week <= :to
group by item
order by week, commits desc`

const repoTopQuery = `
select
    %s as item,
    sum(c.commits) as commits,
    sum(c.additions) - sum(c.deletions) as delta
from contribs c
join members m on
    c.author = m.user and
    c.owner = m.org
join teams t on
    m.team_id = t.id
where
    m.id is not null and
    c.owner = :org and
    c.repo = :repo and
    c.week >= :from and
    c.week <= :to
group by item
order by commits desc`

const repoActivityQuery = `
select
    c.week as weel,
    %s as item,
    sum(c.commits) as commits,
    sum(c.additions) - sum(c.deletions) as delta
from contribs c
join members m on
    c.author = m.user and
    c.owner = m.org
join teams t on
    m.team_id = t.id
where
    m.id is not null and
    c.owner = :org and
    c.repo = :repo and
    c.week >= :from and
    c.week <= :to
group by week, item
order by commits desc`

func StatOrgTop(p map[string]interface{}) (res []StatItem) {
	defer measure("StatOrgTop", time.Now())
	mustSelectN(&res, fmt.Sprintf(orgTopQuery, p["item"], p["sort"]), p)
	return
}

func StatOrgActivity(p map[string]interface{}) (res []StatPoint) {
	defer measure("StatOrgActivity", time.Now())
	mustSelectN(&res, fmt.Sprintf(orgActivityQuery, p["item"], p["sort"]), p)
	return
}

func StatTeamTop(p map[string]interface{}) (res []StatItem) {
	defer measure("StatTeamTop", time.Now())
	mustSelectN(&res, fmt.Sprintf(teamTopQuery, p["item"], p["sort"]), p)
	return
}

func StatTeamActivity(p map[string]interface{}) (res []StatPoint) {
	defer measure("StatTeamActivity", time.Now())
	mustSelectN(&res, fmt.Sprintf(teamActivityQuery, p["item"], p["sort"]), p)
	return
}

func StatUserTop(p interface{}) (res []StatItem) {
	defer measure("StatUserTop", time.Now())
	mustSelectN(&res, userTopQuery, p)
	return
}

func StatUserActivity(p interface{}) (res []StatPoint) {
	defer measure("StatUserActivity", time.Now())
	mustSelectN(&res, userActivityQuery, p)
	return
}

func StatRepoTop(p interface{}) (res []StatItem) {
	defer measure("StatRepoTop", time.Now())
	mustSelectN(&res, repoTopQuery, p)
	return
}

func StatRepoActivity(p interface{}) (res []StatPoint) {
	defer measure("StatRepoActivity", time.Now())
	mustSelectN(&res, repoActivityQuery, p)
	return
}
