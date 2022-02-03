package github

import (
	"github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/api"
	graphql "github.com/cli/shurcooL-graphql"
	"time"
)

// GetOrganizationRepos returns an aggregated list of all repositories
// within a GitHub organization, either private or public
func GetOrganizationRepos(organization string, batchSize int) ([]string, error) {
	opts := api.ClientOptions{
		EnableCache: true,
		Timeout:     10 * time.Second,
	}
	client, err := gh.GQLClient(&opts)
	if err != nil {
		return nil, err
	}
	var orgRepositories = make([]string, 0)
	repositories, totalPagesCount, endCursor, err := getOrganizationRepositories(&client, organization, batchSize)
	for _, repo := range repositories {
		orgRepositories = append(orgRepositories, repo)
	}
	var after = endCursor
	if totalPagesCount > batchSize {
		for {
			repositories, endCursor, err := getOrganizationRepositoriesAfter(&client, organization, after, batchSize)
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

func getOrganizationRepositories(client *api.GQLClient, organization string, batchSize int) ([]string, int, string, error) {
	var query struct {
		Organization struct {
			Repositories struct {
				TotalCount int
				PageInfo   struct {
					StartCursor string
					EndCursor   string
				}
				Nodes []struct {
					Name string
				}
			} `graphql:"repositories(first: $first)"`
		} `graphql:"organization(login: $org)"`
	}

	variables := map[string]interface{}{
		"first": graphql.Int(batchSize),
		"org":   graphql.String(organization),
	}

	err := (*client).Query("OrganizationRepositories", &query, variables)
	if err != nil {
		return nil, 0, "", err
	}

	repositories := make([]string, 0)
	for _, node := range query.Organization.Repositories.Nodes {
		repositories = append(repositories, node.Name)
	}
	return repositories, query.Organization.Repositories.TotalCount, query.Organization.Repositories.PageInfo.EndCursor, nil
}

func getOrganizationRepositoriesAfter(client *api.GQLClient, organization string, after string, batchSize int) ([]string, string, error) {
	var query struct {
		Organization struct {
			Repositories struct {
				TotalCount int
				PageInfo   struct {
					StartCursor string
					EndCursor   string
				}
				Nodes []struct {
					Name string
				}
			} `graphql:"repositories(first: $first, after: $after)"`
		} `graphql:"organization(login: $org)"`
	}

	variables := map[string]interface{}{
		"first": graphql.Int(batchSize),
		"org":   graphql.String(organization),
		"after": graphql.String(after),
	}

	err := (*client).Query("OrganizationRepositories", &query, variables)
	if err != nil {
		return nil, "", err
	}

	repositories := make([]string, 0)
	for _, node := range query.Organization.Repositories.Nodes {
		repositories = append(repositories, node.Name)
	}
	return repositories, query.Organization.Repositories.PageInfo.EndCursor, nil
}
