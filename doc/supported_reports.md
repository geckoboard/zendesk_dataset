# Supported Reports

**Missing a report ?** that you would like to report into Geckoboard. Please mention it on the
Geckoboard community to see if others will find it useful and we will implement it within due time.

## Ticket counts

Ticket counts uses the Zendesk search api documented [here](https://developer.zendesk.com/rest_api/docs/core/search).

We only currently support type:ticket however if there is enough uptake we can support the other
types available in the search api.

The ticket counts can be used to get total counts of anything supported inside the query for
type:ticket. Some examples would be the following;

* Tickets status solved in the last x days/months
* Tickets status x in the last x days/months
* Tickets with tags
* Tickets with tags created in the last x days/months

You can also further group the above for example;

You are using tags to describe tickets for a region, with the group_by config option this is now
possible just by telling us to group by it, and then it will report the count for each tag in the
list.

### Example Report Section

#### Tickets status created in the last 30 days

```yaml
  - name: ticket_counts
    dataset: zendesk.tickets.solved.in.last.month
    filter:
      date_range:
        past: 30
        unit: day
```


#### Tickets status not solved in the last 14 days

```yaml
  - name: ticket_counts
    dataset: zendesk.tickets.not.solved.in.last.fortnight
    filter:
      date_range:
        past: 1
        unit: month
      value:
        'status<': solved
```

#### Tickets matching tag and still open and updated in the last month

```yaml
  - name: ticket_counts
    dataset: zendesk.tickets.pending.reply.still.update.in.last.month
    filter:
      date_range:
        attribute: updated
        past: 1
        unit: month
      value:
        'status:': open
        'tags:': pending_customer_reply
```

#### Tickets grouped by some tags in the last week

```yaml
  - name: ticket_counts
    dataset: zendesk.tickets.created.in.last.week.by.region
    grouped_by:
      key: "tags:"
      name: "Tags"
    filter:
      date_range:
        past: 7
        unit: day
      values:
        'tags:':
        - uk
        - ireland
        - spain
```

## Ticket counts by day

This is one of the newest additions you are now able to use the above filters in ticket counts
but just update the report name to be `ticket_counts_by_day`.  However the group\_by is not
supported.

This only supports grouping by day because Geckoboard Datasets are powerful enough to allow
you to further bucket the data by month or year.

This report template allows you to then plot a line chart because the x-axis is a date field.

### Example Report - Tickets created in the last 6 months

You can use some of the above ticket count examples however remember to change the name
to `ticket_counts_by_day`.

One common report for counts by day would be the following example

```yaml
  - name: ticket_counts_by_day
    dataset: zendesk.tickets.solved.in.last.6.months
    filter:
      date_range:
        past: 6
        unit: month
```


## Detailed ticket metrics

This is another new addition which utilizes two different endpoints. First it uses the search api
to get a filter of the tickets - just like ticket counts filters.

Then it collates all the ticket ids and makes a request to `tickets/show_many.json` sideloading the
metrics which is detailed [here](https://developer.zendesk.com/rest_api/docs/core/side_loading) this
then allows us to retrieve the metrics. Based on a filter just like above.

You are able to report on the following metrics from a ticket metric

* reply\_time
* first\_resolution\_time
* full\_resolution\_time
* agent\_wait\_time
* requester\_wait\_time,
* on\_hold\_time

You can also from each of those metrics choose whether you want to use the **business** or **calendar** value.
Not sure which one to choose read [here](https://support.zendesk.com/hc/en-us/articles/205951808-Calculating-first-reply-time)
for more information on the subject.

Also these metrics are the detailed metrics and allow you to group the metric in time boxes of your
choosing which can then be plotted on a bar/column chart. The **metric options** are required and
all of it sub options.

### Example Reports

#### First reply time using business metric

```yaml
  - name: detailed_metrics
    dataset: zendesk.first.reply.time.last.7.days
    metric_options:
      attribute: reply_time
      unit: business
      grouping:
      - from: 0
        to: 1
        unit: hour
      - from: 1
        to: 8
        unit: hour
      - from: 8
        to: 24
        unit: hour
      - from: 24
        to: 336
        unit: hour
    filter:
      date_range:
        past: 7
        unit: day
```

#### Full resolution time calendar metric which are solved in the last month

```yaml
  - name: detailed_metrics
    dataset: zendesk.full_resolution.reply.time.last.month
    metric_options:
      attribute: full_resolution_time
      unit: calendar
      grouping:
      - from: 0
        to: 8
        unit: hour
      - from: 8
        to: 24
        unit: hour
      - from: 24
        to: 72
        unit: hour
      - from: 72
        to: 772
        unit: hour
    filter:
      date_range:
        past: 1
        unit: month
      value:
        'status:': solved
```
