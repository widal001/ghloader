# GitHub project loader

Bulk update GitHub project items and fields from a TSV or CSV file.

### Quickstart

### Setup

1. Clone the repo
2. Create a [personal access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-personal-access-token-classic) with the `project` permissions
3. Make sure Go in installed on your machine: [Go installation instructions](https://go.dev/doc/install)
4. Run `make setup` and enter your newly created access token when prompted

### Running in web app mode

1. Run `make web-up` to spin up the application
2. Navigate to `http://localhost:8080` to view the batch update form
3. Fill out and submit the form to trigger a batch update of your project
4. Run `make web-down` to tear down the application when you're done

### Running in CLI mode

Run `make cli url=<GitHub project URL> file=<Path to csv or tsv with GitHub project data>` to add or update the items in a project with data from a local TSV or CSV file
