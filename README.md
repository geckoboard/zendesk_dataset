# Zendesk Dataset

Zendesk datasets has been a little side project to help more customers get the data they need from Zendesk into 
[Geckoboard](https://www.geckoboard.com). This uses the latest [Datasets API](https://www.geckoboard.com/whats-new/9)
which is currently in beta but has some great benefits over the older custom widgets, such as one dataset can be used
by multiple widgets as the data isn't tied to a specific widget type.

## What this is and isn't
This **isn't** a replacement for the current [Zendesk integration](https://www.geckoboard.com/integrations/zendesk),
but there are many use cases for reports that are not possible with the current integration and are specific to each
customer. 

So this is where this **will help**.

This program utilizes the zendesk search api and attempts to wrap it up into a config file with a predefined 
Dataset schema for that report template, so if you wanted to get a count of the number of tickets created in the last
30 days matching some tags then it possible to do so with a few options and your api keys for both services.

It is early days and currently only supports getting ticket counts for anything that the 
[Zendesk search api](https://developer.zendesk.com/rest_api/docs/core/search) currently supports

## Missing report that you need?
You can either submit a pull request. Or raise [new issue](https://github.com/geckoboard/zendesk_dataset/issues/new)

Please check if an existing issue is already present before submitting a new one, feel free to comment on an existing
one based on your use case if it very similar.

## Getting Started
As this is Go there won't be the need to install any libraries or anything to get started. 
We have already compiled the code into a distributable binary for the main operating systems/arch select it below

**TODO: Compile binaries for OSX, Linux, Windows, i386, x86_64 and list here to the releases section**

## Usage

We first need to build a configuration file you can start with the following [dummy template](fixtures/example.conf)
which we will describe.

### Authentication
So that we can pull data from Zendesk and push data to Geckoboard we need to first get login api keys to support this

##### Geckoboard 
In the below example you only need to change the `api_key` which you can get from your account section.
The url in this example is correct and should never change.

##### Zendesk
We are supporting two authentication options Password and Apikey access. 
To use email/password auth only supply `email` and `password options` if want to authenticate with an apikey.

First generate one from Admin > Channels > API in Zendesk and supply only the `api_key` and `email` options.

In both cases a subdomain must be supplied.

```json
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
```

### Reports (the fun part)

So reports sit under the Zendesk key in the json configuration file. 
Again using the example.conf this is one possible report utilizing both the status and tags, and a second report just
specifying all tickets in the last 3 months.

Lets go into some more detail about each option below;

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

The `name` is the template specified by this application at the moment we only support the report template called `ticket_counts`

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

## Running the program

Now that we have a configuration file we should be ready to run it. With the downloaded binary ensure you are in the
same directory as the binary and run the following;

```sh
./zendesk_datasets -config your_config_file_full_path
```

While the program is running you should see output like the following as an example;
If there is an error it will output the error that occurred, otherwise will tell you all was successful.

```
$ ERRO: Processing report 'your.report.1' failed with: Custom input requires the operator one of [< : >]
$ INFO: Processing report 'your.report.2' completed successfully
$ Completed processing all reports...
```


## Example
The below config example will get you up and running to get your first dataset created from zendesk ticket counts created
in the last 30 days. Try it all you need to do is replacing the geckoboard api key and your zendesk credentials.

```
{
    "geckoboard": {
        "api_key": "Ap1K4y",
        "url": "https://testing.geckoboardexample.com"
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
