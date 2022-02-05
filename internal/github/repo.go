package github

import (
	"fmt"
	"github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/api"
	graphql "github.com/cli/shurcooL-graphql"
	"log"
	"strings"
	"time"
)

// GetOrganizationRepos returns an aggregated list of all repositories
// within a GitHub organization, either private or public
func GetOrganizationRepos(organization string, query string, batchSize int) ([]string, error) {
	organizationSearchQuery := fmt.Sprintf("org:\"%s\"", organization)
	var queryString string
	if strings.Contains(query, organizationSearchQuery) {
		queryString = query
	} else {
		if query != "" {
			queryString = fmt.Sprintf("%s %s", organizationSearchQuery, query)
		} else {
			queryString = organizationSearchQuery
		}
	}
	if !strings.Contains(query, "fork:") {
		//Include forks by default
		queryString += " fork:true"
	}
	log.Println("[debug] queryString", queryString)

	opts := api.ClientOptions{
		EnableCache: true,
		Timeout:     10 * time.Second,
	}
	client, err := gh.GQLClient(&opts)
	if err != nil {
		return nil, err
	}
	var orgRepositories = make([]string, 0)
	repositories, repositoryCount, endCursor, err := getOrganizationRepositories(&client, queryString, batchSize)
	for _, repo := range repositories {
		orgRepositories = append(orgRepositories, repo)
	}
	var after = endCursor
	if repositoryCount > batchSize {
		for {
			repositories, endCursor, err := getOrganizationRepositoriesAfter(&client, queryString, batchSize, after)
			if err != nil {
				return nil, err
			}
			for _, repo := range repositories {
				orgRepositories = append(orgRepositories, repo)
			}
			if endCursor == "" {
				break
			}
			after = endCursor
		}
	}
	return orgRepositories, nil
}

type RepositoryFragment struct {
	Name string
}

func getOrganizationRepositories(client *api.GQLClient, queryString string, batchSize int) ([]string, int, string, error) {
	/*
		search(type: REPOSITORY, query: $query, first: $first) {
		    pageInfo {
		      endCursor
		      startCursor
		    }
		    repositoryCount
		    repos: edges {
		      repo: node {
		        ... on Repository {
		          name
			  }
			}
		}
	*/
	var query struct {
		Search struct {
			PageInfo struct {
				StartCursor string
				EndCursor   string
			}
			RepositoryCount int
			Repos           []struct {
				Repo struct {
					RepositoryFragment `graphql:"... on Repository"`
				} `graphql:"repo: node"`
			} `graphql:"repos: edges"`
		} `graphql:"search(type: REPOSITORY, query: $query, first: $first)"`
	}
	variables := map[string]interface{}{
		"query": graphql.String(queryString),
		"first": graphql.Int(batchSize),
	}

	err := (*client).Query("OrganizationRepositories", &query, variables)
	if err != nil {
		return nil, 0, "", err
	}

	repositories := make([]string, 0)
	for _, r := range query.Search.Repos {
		repositories = append(repositories, r.Repo.Name)
	}
	return repositories, query.Search.RepositoryCount, query.Search.PageInfo.EndCursor, nil
}

func getOrganizationRepositoriesAfter(client *api.GQLClient, queryString string, batchSize int, after string) ([]string, string, error) {
	var query struct {
		Search struct {
			PageInfo struct {
				StartCursor string
				EndCursor   string
			}
			RepositoryCount int
			Repos           []struct {
				Repo struct {
					RepositoryFragment `graphql:"... on Repository"`
				} `graphql:"repo: node"`
			} `graphql:"repos: edges"`
		} `graphql:"search(type: REPOSITORY, query: $query, first: $first, after: $after)"`
	}
	variables := map[string]interface{}{
		"query": graphql.String(queryString),
		"first": graphql.Int(batchSize),
		"after": graphql.String(after),
	}

	err := (*client).Query("OrganizationRepositoriesAfter", &query, variables)
	if err != nil {
		return nil, "", err
	}

	repositories := make([]string, 0)
	for _, r := range query.Search.Repos {
		repositories = append(repositories, r.Repo.Name)
	}
	return repositories, query.Search.PageInfo.EndCursor, nil
}
