# GitHub project loader

Bulk update GitHub project items and fields from a TSV or CSV file.

### Quickstart

1. Clone the repo
2. Create a [personal access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-personal-access-token-classic) with the `project` permissions
3. Run `make setup` and enter your newly created access token when prompted
4. Run `make run url=<GitHub project URL> file=<Path to csv or tsv with GitHub project data>` to add or update the items in a project with data from a local TSV or CSV file
