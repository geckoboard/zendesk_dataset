# Modifying the report (the fun part!)

To modify the data pulled back from Zendesk you will need to edit the reports section of the JSON configuration file. Multiple reports can be set up in the same configuration file.

In the example below the first report pulls back the number of open tickets tagged `beta` and `freetrial`, created in the past 14 days, grouped by tag. The second report pulls back the total number of tickets created in the last 3 months.

Lets go into some more detail about each option below:

```yaml
reports:
- name: ticket_counts
  dataset: your.report.1
  group_by:
    key: 'tags:'
    name: Tags
  filter:
    date_range:
    - past: 14
      unit: day
    value:
      'status:': open
    values:
      'tags:':
      - beta
      - freetrial
- Name: ticket_counts
  dataset: your.report.2
  filter:
    date_range:
    - past: 3
      unit: month

```
### Options

#### Name

The `name` is the name of the report template to be used. Checkout what reports are supported and some examples [here](supported_reports.md). As we expand the application, more templates will become available.

```yaml
name: ticket_counts,
```

#### Dataset

The `dataset` is where you specify the name of the dataset that will be created in Geckoboard.

```yaml
dataset: your.report.1
```

#### Filter

The `filter` option is where the search filter for Zendesk is specified.

#### Date range

The `date_range` option allows you to specify a date range. In the example below, we're asking for
tickets created in the the last 14 days (calculated from today's date).

By default we assume you want to use the `created` date. If this is not the case you can specify an
additional key called `attribute` which allows one of the following options: `["created", "updated", "solved", "due_date"]`.

```yaml
date_range:
 - past: 14
   unit: day
```

There is also support for hardcoded dates if you want finer control over the date range.
For example, to get the count of solved tickets in first quarter of 2016 you would use the options shown in the config below.

Notice that in the custom value we have used the operators `>=` and `<`. You must specify an operator as part of the `custom` attribute, and it can be any of the ones outline [here](https://support.zendesk.com/hc/en-us/articles/203663226#topic_ngr_frb_vc).

```yaml
date_range:
- attribute: solved
  custom: ">=2016-01-01"
- attribute: solved
  custom: "<2016-04-01"
```
#### Value

The `value` option allows you to specify an attribute supported by the Zendesk Search API and its value. Note that the key **must**
again include an operator. In the example below we ask for `status: open`, meaning tickets matching the status exactly open.
By contrast, `"status<": "solved"` will return all unsolved tickets.

```yaml
value:
  'status:': open
```

#### Values

The `values` option is similar to the `value` option in that you need to specify the operator in the key, however
you specify an array of the tags just once. Using values also allows you to group by that key

```yaml
values:
  'tags:':
   - beta
   - freetrial
```

#### Group by

The `group_by` option allows you to group by a key in the `values` object. `group_by` is an object with two keys - `key` and `name`. The `key` must exactly match one of the keys in the `values` object or it will error.

The `name` attribute sets how this grouping will be displayed in the Geckoboard dataset. You can omit the name and it will use the key instead.

Note that `group_by` doesn't sit under the `filter` key, but rather the `report` key.

This is an example of you'd `group_by` the tags we defined in the previous example:

```yaml
group_by:
  key: 'tags:'
  name: Tags
```

This example would return the counts for both tags:beta and tags:freetrial seperately from each other, but the results be combined with all the other filters specified.
