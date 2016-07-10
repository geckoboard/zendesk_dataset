# Getting Started

## 1. Download the correct Binary

As this is a Go application there's no need to install any libraries to get started.

Just click on the distributable binary that matches your operating system and architecture below to download it:

**TODO: Compile binaries for OSX, Linux, Windows, i386, x86_64 and link here to the releases section**

## 2. Build the configuration file

The configuration file is where you enter the your Geckoboard and Zendesk API keys, and where you specify what data to pull from Zendesk. Get started by copying the example below and saving it as a JSON file.

When we run our application with this configuration file it will create a Dataset in Geckoboard called `tickets.created.in.last.30.days` that contains - unsurprisingly - the number of tickets created in Zendesk in the last 30 days!

### Modifying the configuration file

Before you can run the application you will need to modify the config file.

First you will need to edit the Geckoboard `api_key` to match the one found in the Account section of your Geckoboard account. You won't need to edit the `url`.

You can authenticate with Zendesk using either a password or an API key. To authenticate with email and password supply the `email` and `password` options. To authenticate with an API key you'll first generate one in Zendesk by heading to Admin > Channels > API. Then, in the config file, supply the `api_key` and `email` options. **In both cases your Zendesk `subdomain` must be supplied.**

```
{
    "geckoboard": {
        "api_key": "Ap1K4y",
        "url": "https://api.geckoboard.com"
    },
    "zendesk": {
        "auth": {
            "api_key": "12345",
            "email": "test@example.com",
            "password": "test",
            "subdomain": "testing"
        },
        "reports": [
            {
                "Name": "ticket_counts",
                "dataset": "tickets.created.in.last.30.days",
                "filter": {
                    "date_range": [
                        {
                            "attribute": "created",
                            "past": 30,
                            "unit": "day"
                        }
                    ]
                }
            }
        ]
    }
}
```

## 3. Run the program

Now that we have a configuration file we're ready to run it. In the terminal ensure you are in the
same directory as the binary and run the following:

```sh
./zendesk_datasets -config full_path_to_your_config_file
```

While the program is running you should see output like the following:

```
$ ERRO: Processing report 'your.report.1' failed with: Custom input requires the operator one of [< : >]
$ INFO: Processing report 'your.report.2' completed successfully
$ Completed processing all reports...
```

If an error occurs, it'll be output to the console. Otherwise, you'll be told all was successful!

## 4. Building a widget from the Dataset

Head to Geckoboard, click 'Add Widget', and select the Datasets integration. In the pop-out panel that appears you should see your new dataset `tickets.created.in.last.30.days`. You can use this to build a widget showing your Zendesk ticket count.

