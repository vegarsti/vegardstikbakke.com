---
title: "Extract Table API"
date: 2024-08-17
slug: extract-table
---

How great would it be to pull a table from an image or PDF and instantly turn it into structured text?
![The table](https://results.extract-table.com/showcase.png)

I built this! Check out [extract-table.com](https://extract-table.com).

![The table](https://results.extract-table.com/screenshot.png)

Uploading an image produces a result like in the screenshot above, hosted at [result.extract-table.com](https://results.extract-table.com/52bc7d95f07a30eef023e36bde0de9c3aab7e057aaa854b44419291ec56b0839).

Please try it!

---

The site uses `api.extract-table.com` which can also be used directly.

```
$ curl https://api.extract-table.com -X \
    POST -H "Content-Type: image/png" \
    -H "Accept: text/csv" \
    --data-binary @invoice.png
Items,Type,,Quantity,Price,Total
Test lonut,Rate,Plan,2,$11.00/$242.00,$253.00
bulk import item2,MRC,,1,$2.20,$2.20
bulk import item2 desc,,,,,
bulk import item4,NRC,,2,$4.40,$8.80
bulk import item4 desc,,,,,
bulk import item7,NRC,,2,$0.00,$0.00
bulk import item7 desc,,,,,
asdasd,NRC,,2,$0.00,$0.00
asdasd,,,,,
mine,MRC,,2,$121.00,$242.00
2bulk import item1,MRC,,1,$33.00,$33.00
bulk import item 1 desc,,,,,
```

The API is an AWS Lambda function which uses [Amazon Textract](https://aws.amazon.com/textract/), AWS's OCR service that gives you all words with their bounding box.
It even has a mode for detecting tables, which is 10x more expensive.

I actually built this in 2021 and [posted about it on Hacker News](https://news.ycombinator.com/item?id=28680136).
Since then the API has been used 26 thousand times.
I originally used the table detection algorithm, but that became too expensive to host for free.

I've since written my own table detection algorithm, which was used to produce the result above.
But explaining that algorithm is for another post!

[Browse the code](https://github.com/vegarsti/extract-table).
