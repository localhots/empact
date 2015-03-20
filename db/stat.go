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

const joinContribFT = `
join team_members tm on
    c.user_id = tm.user_id and
    c.org_id = tm.org_id
join orgs o on
    c.org_id = o.id
join teams t on
    tm.team_id = t.id
join users u on
    c.user_id = u.id
join repos r on
    c.repo_id = r.id`

func StatOrgTop(p map[string]interface{}) (res []StatItem) {
	defer measure("StatOrgTop", time.Now())
	mustSelectN(&res, fmt.Sprintf(`
        select
            %s as item,
            sum(c.commits) as commits,
            sum(c.additions) - sum(c.deletions) as delta
        from contribs c
        `+joinContribFT+`
        where
            o.login = :org and
            c.week >= :from and
            c.week <= :to
        group by item
        order by commits desc
    `, p["item"]), p)
	return
}

func StatOrgActivity(p map[string]interface{}) (res []StatPoint) {
	defer measure("StatOrgActivity", time.Now())
	mustSelectN(&res, fmt.Sprintf(`
        select
            %s as item,
            sum(c.commits) as commits,
            sum(c.additions) - sum(c.deletions) as delta,
            c.week as week
        from contribs c
        `+joinContribFT+`
        where
            o.login = :org and
            c.week >= :from and
            c.week <= :to
        group by item, week
        order by week, commits desc
    `, p["item"]), p)
	return
}

func StatTeamTop(p map[string]interface{}) (res []StatItem) {
	defer measure("StatTeamTop", time.Now())
	mustSelectN(&res, fmt.Sprintf(`
        select
            %s as item,
            sum(c.commits) as commits,
            sum(c.additions) - sum(c.deletions) as delta
        from contribs c
        `+joinContribFT+`
        where
            o.login = :org and
            t.name = :team and
            c.week >= :from and
            c.week <= :to
        group by item
        order by commits desc
    `, p["item"]), p)
	return
}

func StatTeamActivity(p map[string]interface{}) (res []StatPoint) {
	defer measure("StatTeamActivity", time.Now())
	mustSelectN(&res, fmt.Sprintf(`
        select
            %s as item,
            sum(c.commits) as commits,
            sum(c.additions) - sum(c.deletions) as delta,
            c.week as week
        from contribs c
        `+joinContribFT+`
        where
            o.login = :org and
            t.name = :team and
            c.week >= :from and
            c.week <= :to
        group by item, week
        order by week, commits desc
    `, p["item"]), p)
	return
}

func StatUserTop(p map[string]interface{}) (res []StatItem) {
	defer measure("StatUserTop", time.Now())
	mustSelectN(&res, `
        select
            r.name as item,
            sum(c.commits) as commits,
            sum(c.additions) - sum(c.deletions) as delta
        from contribs c
        join orgs o on
            c.org_id = o.id
        join users u on
            c.user_id = u.id
        join repos r on
            c.repo_id = r.id
        where
            o.login = :org and
            u.login = :user and
            c.week >= :from and
            c.week <= :to
        group by item
        order by commits desc
    `, p)
	return
}

func StatUserActivity(p map[string]interface{}) (res []StatPoint) {
	defer measure("StatUserActivity", time.Now())
	mustSelectN(&res, `
        select
            c.week as week,
            r.name as item,
            sum(c.commits) as commits,
            sum(c.additions) - sum(c.deletions) as delta
        from contribs c
        join orgs o on
            c.org_id = o.id
        join users u on
            c.user_id = u.id
        join repos r on
            c.repo_id = r.id
        where
            o.login = :org and
            u.login = :user and
            c.week >= :from and
            c.week <= :to
        group by item
        order by week, commits desc
    `, p)
	return
}

func StatRepoTop(p map[string]interface{}) (res []StatItem) {
	defer measure("StatRepoTop", time.Now())
	mustSelectN(&res, fmt.Sprintf(`
        select
            %s as item,
            sum(c.commits) as commits,
            sum(c.additions) - sum(c.deletions) as delta
        from contribs c
        `+joinContribFT+`
        where
            o.login = :org and
            r.name = :repo and
            c.week >= :from and
            c.week <= :to
        group by item
        order by commits desc
    `, p["item"]), p)
	return
}

func StatRepoActivity(p map[string]interface{}) (res []StatPoint) {
	defer measure("StatRepoActivity", time.Now())
	mustSelectN(&res, fmt.Sprintf(`
        select
            c.week as week,
            %s as item,
            sum(c.commits) as commits,
            sum(c.additions) - sum(c.deletions) as delta
        from contribs c
        `+joinContribFT+`
        where
            o.login = :org and
            r.name = :repo and
            c.week >= :from and
            c.week <= :to
        group by week, item
        order by commits desc
    `, p["item"]), p)
	return
}
