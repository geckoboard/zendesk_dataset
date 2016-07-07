# Zendesk Datasets

## Background

Zendesk datasets is a small Go application that makes use of [Geckoboard's](https://www.geckoboard.com) new [Datasets API](https://www.geckoboard.com/whats-new/9)
(currently in Beta) to pull data from Zendesk into Geckoboard.

The new Datasets API has major advantages over Geckoboard's older Custom Widgets API including the ability to power multiple widgets from one dataset and the ability to switch visualisaton without making code changes. It also now includes powerful grouping and time bucketing features.   

While this application **isn't** intended as a replacement for the Geckoboard's built-in [Zendesk integration](https://www.geckoboard.com/integrations/zendesk) it does make it possible to build some reports that can't currently be achieved. 

Currently it supports getting ticket counts for anything permitted by the 
[Zendesk search api](https://developer.zendesk.com/rest_api/docs/core/search), but we plan to extend the functionality soon.

If you have any feedback we're very interested to hear it. You can either submit a pull request. Or raise [new issue](https://github.com/geckoboard/zendesk_dataset/issues/new)


## Getting Started

### 1. Download the correct Binary

As this is a Go application there's no need to install any libraries to get started.

Just click on the distributable binary that matches your operating systems/arch below to download it:

**TODO: Compile binaries for OSX, Linux, Windows, i386, x86_64 and list here to the releases section**

### 2. Build the configuration file

We now need to build a configuration file for the application. The configration file is where you enter the your Geckoboard and Zendesk API keys, and how you configure what data to pull from Zendesk. To get started copy the example below and save it as a JSON file.

When we run our application with configuration file it will create a Dataset in Geckoboard called `tickets.created.in.last.30.days` that pulls the number of tickets created in Zendesk in the last 30 days.

#### Modifying the configuration file

Before you can run the application you will need to modify the config file. 

First you will need to edit the Geckoboard `api_key` to match the one found in the Account section of your Geckoboard account. You should not need to edit the `url`.

With Zendesk you can either authenticate with  Password or API key. To use email/password auth only supply `email` and `password` options. To authenticate with the API first generate one from Admin > Channels > API in Zendesk and then supply only the `api_key` and `email` options. **In both cases your Zendesk `subdomain` must be supplied.**

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
./zendesk_datasets -config your_config_file_full_path
```

While the program is running you should see output like the following:

```
$ ERRO: Processing report 'your.report.1' failed with: Custom input requires the operator one of [< : >]
$ INFO: Processing report 'your.report.2' completed successfully
$ Completed processing all reports...
```

If there is an error it will output the error that occurred, otherwise will tell you all was successful.

### 4. Building a widget from the Dataset

Now if you go to add a widget in Geckoboard and select Datasets, in the picker you should see `tickets.created.in.last.30.days`. You can now use this to build a widget showing your Zendesk ticket count. 


## Modifying the report (the fun part!)

To modify the data pulled back from Zendesk you will need to edit the reports section of the JSON configuration file. Multiple reports can be created in the same configuration file.

In the example below the first report pulls back the number of open tickets tagged `beta` and `freetrial` created in the past 14 days, grouped by Tag. The second report meanwhile just pulls back the number of tickets created in the past 3 months.

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

The `name` is the template specified by this application at the moment we only support the report template called `ticket_counts`. As we expand the application, more templates will become available.

```json
"Name": "ticket_counts",
```

##### Dataset

The `dataset` is the dataset name that will be created in Geckoboard it following the naming scheme of the following example

```json
"dataset": "your.report.1"
```

##### Filter

The `filter` option is where the search filter for Zendesk is specified inside that we have the following;

##### Date range

The `date_range` option allows you to specify a date range in the below example ask for
tickets created in the the last 14 days which is calculated from today date. 

By default we assume you want to use the created date if this is not the case you can specify an 
additional key called `attribute` which allows one of the following ["created", "updated", "solved", "due_date"] just like
the search api supports.

```json
"date_range": [
    {
        "past": 14,
        "unit": "day"
    }
]
```

There is also support for hardcoded dates should you want to be specific on the date range. 
For example to get the count of solved tickets in first quarter of 2016 you can specify the below config.

Notice that in the custom value we have some operators specified `>=` `<` these are required for the query to be valid
on the zendesk request. It can be one of the following operators outlined [here](https://support.zendesk.com/hc/en-us/articles/203663226#topic_ngr_frb_vc)

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

The `value` option allows you to specify an attribute supported by the search api and its value. Note that the key **must**
again include the operator in this example we ask for status: open - meaning tickets matching the status exactly open.
You could have `"status<": "solved"` which will return all unsolved tickets.

```json
"value": {
  "status:": "open"
}
```

##### Values

The `values` option is really similar to the `value` option in that you need to specify the operator in the key, however
you just specify an array of the tags just once. Using values also allows you to group by that key 

```json
"values": {
    "tags:": [
        "beta",
        "freetrial"
    ]
}
```

##### Group by

The `group_by` option allows you to group by a values key. Using the example just above this one we can specify to
group by the key must exactly as the values key or it will error so it would be like below.

The name attribute is what to display that grouping in the geckoboard dataset as you probably don't want tags: as
the grouping value, but you can omit the name and it will use the key instead.

Note that group_by doesn't sit under filter key but rather the report key

```json
"group_by": {
  "key": "tags:",
  "name": "Tags"
}
```

and in the above example it would return the counts both tags:beta and tags:freetrial seperately from eachother
but will be combined with all the other filters specified.
