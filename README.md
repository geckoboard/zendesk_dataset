# Zendesk Datasets

## Start here...

* [Getting Started](doc/getting_started.md)
* [Supported Reports](doc/supported_reports.md)
* [Modifying the reports](doc/modifying_report.md)


## Background

Zendesk Datasets is a small Go application that makes use of [Geckoboard's](https://www.geckoboard.com) new [Datasets API](https://www.geckoboard.com/whats-new/9)
(currently in Beta) to pull data from Zendesk into Geckoboard.

The new Datasets API has major advantages over Geckoboard's older Custom Widgets API including the ability to power multiple widgets from one dataset and the ability to switch visualisaton without making code changes. It also includes powerful grouping and time bucketing features.

While this application **isn't** intended as a replacement for the Geckoboard's built-in [Zendesk integration](https://www.geckoboard.com/integrations/zendesk) it does make it possible to build some widgets that the built-in one doesn't offer.

Zendesk Datasets currently supports getting ticket counts for anything permitted by the
[Zendesk search api](https://developer.zendesk.com/rest_api/docs/core/search), but we plan to extend the functionality soon.

We're very interested to hear any feedback you might have. You can either submit a pull request or raise a [new issue](https://github.com/geckoboard/zendesk_dataset/issues/new).

## YAML config supported

We now support yaml configuration format over json due to it readability and cleaness. We do however
support the json format if you have already started using the program.

However if you would like to transform your current json config into the yaml format you can do by
using a json to yaml convert online like `http://www.json2yaml.com/` **please remember to remove
any api keys before doing so as a precaution of safety.**



