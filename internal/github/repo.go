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

type RepositoryInfo struct {
	Name    string
	IsEmpty bool
}

type organizationResponse struct {
	repositories    []RepositoryInfo
	repositoryCount int
	endCursor       string
}

// GetOrganizationRepos returns an aggregated list of all repositories
// within a GitHub organization, either private or public
func GetOrganizationRepos(organization string, query string, batchSize int) ([]RepositoryInfo, error) {
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
	var orgRepositories = make([]RepositoryInfo, 0)
	organizationResponse, err := getOrganizationRepositories(&client, queryString, batchSize)
	if err != nil {
		return nil, err
	}
	orgRepositories = append(orgRepositories, organizationResponse.repositories...)
	var after = organizationResponse.endCursor
	if organizationResponse.repositoryCount > batchSize {
		for {
			organizationResponse, err := getOrganizationRepositoriesAfter(&client, queryString, batchSize, after)
			if err != nil {
				return nil, err
			}
			orgRepositories = append(orgRepositories, organizationResponse.repositories...)
			if organizationResponse.endCursor == "" {
				break
			}
			after = organizationResponse.endCursor
		}
	}
	return orgRepositories, nil
}

type RepositoryFragment struct {
	Name    string
	IsEmpty bool
}

func getOrganizationRepositories(client *api.GQLClient, queryString string, batchSize int) (organizationResponse, error) {
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
		              isEmpty
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
		return organizationResponse{}, err
	}

	repositories := make([]RepositoryInfo, 0)
	for _, r := range query.Search.Repos {
		repositories = append(repositories, RepositoryInfo{
			Name:    r.Repo.Name,
			IsEmpty: r.Repo.IsEmpty,
		})
	}
	return organizationResponse{
		repositories:    repositories,
		repositoryCount: query.Search.RepositoryCount,
		endCursor:       query.Search.PageInfo.EndCursor,
	}, nil
}

func getOrganizationRepositoriesAfter(client *api.GQLClient, queryString string, batchSize int, after string) (organizationResponse, error) {
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
		return organizationResponse{}, err
	}

	repositories := make([]RepositoryInfo, 0)
	for _, r := range query.Search.Repos {
		repositories = append(repositories, RepositoryInfo{
			Name:    r.Repo.Name,
			IsEmpty: r.Repo.IsEmpty,
		})
	}
	return organizationResponse{
		repositories:    repositories,
		repositoryCount: query.Search.RepositoryCount,
		endCursor:       query.Search.PageInfo.EndCursor,
	}, nil
}
