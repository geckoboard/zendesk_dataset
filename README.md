# Zendesk Datasets

## Background

Zendesk Datasets is a small Go application that makes use of [Geckoboard's](https://www.geckoboard.com) new [Datasets API](https://www.geckoboard.com/whats-new/9)
(currently in Beta) to pull data from Zendesk into Geckoboard.

The new Datasets API has major advantages over Geckoboard's older Custom Widgets API including the ability to power multiple widgets from one dataset and the ability to switch visualisaton without making code changes. It also includes powerful grouping and time bucketing features.   

While this application **isn't** intended as a replacement for the Geckoboard's built-in [Zendesk integration](https://www.geckoboard.com/integrations/zendesk) it does make it possible to build some widgets that the built-in one doesn't offer.

Zendesk Datasets currently supports getting ticket counts for anything permitted by the 
[Zendesk search api](https://developer.zendesk.com/rest_api/docs/core/search), but we plan to extend the functionality soon.

We're very interested to hear any feedback you might have. You can either submit a pull request or raise a [new issue](https://github.com/geckoboard/zendesk_dataset/issues/new).


## Getting Started

### 1. Download the correct Binary

As this is a Go application there's no need to install any libraries to get started.

Just click on the distributable binary that matches your operating system and architecture below to download it:

**TODO: Compile binaries for OSX, Linux, Windows, i386, x86_64 and link here to the releases section**

### 2. Build the configuration file

The configuration file is where you enter the your Geckoboard and Zendesk API keys, and where you specify what data to pull from Zendesk. Get started by copying the example below and saving it as a JSON file.

When we run our application with this configuration file it will create a Dataset in Geckoboard called `tickets.created.in.last.30.days` that contains - unsurprisingly - the number of tickets created in Zendesk in the last 30 days!

#### Modifying the configuration file

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

### 3. Run the program

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

### 4. Building a widget from the Dataset

Head to Geckoboard, click 'Add Widget', and select the Datasets integration. In the pop-out panel that appears you should see your new dataset `tickets.created.in.last.30.days`. You can use this to build a widget showing your Zendesk ticket count.


## Modifying the report (the fun part!)

To modify the data pulled back from Zendesk you will need to edit the reports section of the JSON configuration file. Multiple reports can be set up in the same configuration file.

In the example below the first report pulls back the number of open tickets tagged `beta` and `freetrial`, created in the past 14 days, grouped by tag. The second report pulls back the total number of tickets created in the last 3 months.

Lets go into some more detail about each option below:

```json
"reports": [
            {
                "Name": "ticket_counts",
                "dataset": "your.report.1",
                "group_by": {
                  "key": "tags:",
                  "name": "Tags"
                },
                "filter": {
                    "date_range": [
                        {
                            "past": 14,
                            "unit": "day"
                        }
                    ],
                    "value": {
                        "status:": "open"
                    },
                    "values": {
                        "tags:": [
                            "beta",
                            "freetrial"
                        ]
                    }
                }
            },
            {
                "Name": "ticket_counts",
                "dataset": "your.report.2",
                "filter": {
                    "date_range": [
                        {
                            "past": 3,
                            "unit": "month"
                        }
                    ]
                }
            }
        ]
```
#### Options

##### Name

The `name` is the name of the report template to be used. At the moment we only support one report template, called `ticket_counts`. As we expand the application, more templates will become available.

```json
"Name": "ticket_counts",
```

##### Dataset

The `dataset` is where you specify the name of the dataset that will be created in Geckoboard.

```json
"dataset": "your.report.1"
```

##### Filter

The `filter` option is where the search filter for Zendesk is specified.

##### Date range

The `date_range` option allows you to specify a date range. In the example below, we're asking for
tickets created in the the last 14 days (calculated from today's date). 

By default we assume you want to use the `created` date. If this is not the case you can specify an 
additional key called `attribute` which allows one of the following options: `["created", "updated", "solved", "due_date"]`.

```json
"date_range": [
    {
        "past": 14,
        "unit": "day"
    }
]
```

There is also support for hardcoded dates if you want finer control over the date range. 
For example, to get the count of solved tickets in first quarter of 2016 you would use the options shown in the config below.

Notice that in the custom value we have used the operators `>=` and `<`. You must specify an operator as part of the `custom` attribute, and it can be any of the ones outline [here](https://support.zendesk.com/hc/en-us/articles/203663226#topic_ngr_frb_vc).

```json
"date_range": [
    {
        "attribute": "solved",
        "custom": ">=2016-01-01"
    },
    {
        "attribute": "solved",
        "custom": "<2016-04-01"
    }
]
```
##### Value

The `value` option allows you to specify an attribute supported by the Zendesk Search API and its value. Note that the key **must**
again include an operator. In the example below we ask for `status: open`, meaning tickets matching the status exactly open.
By contrast, `"status<": "solved"` will return all unsolved tickets.

```json
"value": {
  "status:": "open"
}
```

##### Values

The `values` option is similar to the `value` option in that you need to specify the operator in the key, however
you specify an array of the tags just once. Using values also allows you to group by that key 

```json
"values": {
    "tags:": [
        "beta",
        "freetrial"
    ]
}
```

##### Group by

The `group_by` option allows you to group by a key in the `values` object. `group_by` is an object with two keys - `key` and `name`. The `key` must exactly match one of the keys in the `values` object or it will error.

The `name` attribute sets how this grouping will be displayed in the Geckoboard dataset. You can omit the name and it will use the key instead.

Note that `group_by` doesn't sit under the `filter` key, but rather the `report` key.

This is an example of you'd `group_by` the tags we defined in the previous example:

```json
"group_by": {
  "key": "tags:",
  "name": "Tags"
}
```

This example would return the counts for both tags:beta and tags:freetrial seperately from each other, but the results be combined with all the other filters specified.
