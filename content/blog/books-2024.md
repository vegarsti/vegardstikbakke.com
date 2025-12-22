---
title: Books I read in 2024
slug: books-2024
date: 2025-01-12
---

I finished 9 books this year, in total [2006 pages](https://www.goodreads.com/user/year_in_books/2024/3400170). That is on the lower end of reading since I started tracking, not that anyone's counting!
`▁▃█▅▃▄▃▁`

As usual, recommendations in bold. Below I'll briefly say something about each book.

- **Tidy First? A Personal Exercise in Empirical Software Design** (2024) by Kent Beck
- **The Death of Ivan Ilyich** (1886) by Leo Tolstoy
- **Blood Meridian** (1985) by Cormac McCarthy (_reread_)
- The Spy Who Came In From the Cold (1963) by John le Carré
- As I Lay Dying (1930) by William Faulkner
- **The Undoing Project: A Friendship That Changed Our Minds** (2016) by Michael Lewis
- **The Old Man and the Sea** (1952) by Ernest Hemingway
- Code Health Guardian: The Old-New Role of a Human Programmer in the AI Era (2024) by Artie Shevchenko
- **21 Church Fathers** (2000) by Peter Halldorf

There were a handful of books I spent a meaningful amount of time reading but didn't complete.

- **The Making of the Atomic Bomb** (1986) by Richard Rhodes
- Softwar: An Intimate Portrait of Larry Ellison and Oracle (2003) by Matthew Symonds
- Rust in Action (2021) by Tim McNamara
- Hyperion (1989) by Dan Simmons
- Writing for Developers: Blogs that Get Read by Piotr Sarna and Cynthia Dunlop (2024)

### Fiction

I read more fiction books than non-fiction books this year! 🎉

- The Death of Ivan Ilyich (1886) by Leo Tolstoy. A gripping novella about death, family, career, facades, hypocrisy, and the meaning of life. You should read it!
- Blood Meridian (1985) by Cormac McCarthy. So beautiful! McCarthy's prose is just incredible. See [(Probably not) All of the similes in Blood Meridian](https://biblioklept.org/2017/09/18/probably-not-all-of-the-similes-in-blood-meridian/). And so dark! It's basically about violence and evil. Set in the mid-1800s, follows the Kid who joins a scalp-hunting expedition led by the enigmatic Judge. The Judge is rendered so well that I'm not sure I would watch a movie adaptation if there ever is one.
- The Spy Who Came In From the Cold (1963) by John le Carré. I wanted to try something by the renowned spy novel master, and was recommended this one. It was okay 🤷‍♂️
- As I Lay Dying (1930) by William Faulkner. Set in the South. The mother of a family is dying and each chapter is narrated by one of the grieving family members. It was interesting and engaging, and probably a classic for good reason. There's often talk of Faulkerian writing, but I didn't see anything super particular in his style. The book hasn't been much on my mind since I read it, but worth reading anyway.
- The Old Man and the Sea (1952) by Ernest Hemingway. I read this in two sittings! It's short, exciting and touching.
- Hyperion (1989) by Dan Simmons. As with so many science fiction books, I loved the exposition. But the book is too lengthy and too dialogue heavy, so I abandoned it after 200 pages.

### Non-fiction

- Tidy First? A Personal Exercise in Empirical Software Design (2024) by Kent Beck. A small book (~100 pages) by [Kent Beck](https://en.wikipedia.org/wiki/Kent_Beck) of agile and TDD fame. Very pragmatic and enjoyable book on how to think about the balance between adding new features to a codebase and improving the code.
- The Undoing Project: A Friendship That Changed Our Minds (2016) by Michael Lewis. Michael Lewis needs no introduction. This is about Daniel Kahneman and Amos Tversky, the Israeli psychologists who founded behavioral economics. Super engaging as always with Lewis, and interesting subject.
- Code Health Guardian: The Old-New Role of a Human Programmer in the AI Era (2024) by Artie Shevchenko. Very 2024 book title, but the book is less about AI and more some good recommendations about maintaining a codebase.
- 21 Church Fathers (2000) by Peter Halldorf. Great primer on the most influential people in the Church in the first millennium.
- The Making of the Atomic Bomb (1986) by Richard Rhodes. At almost 900 pages, this is an undertaking! Rhodes is a historian but also a journalist, which shines through because this is a fun read. It covers the people and discoveries that led to nuclear fission and the atomic bomb. Not surprisingly it won the Pulitzer and is [praised in general](https://en.wikipedia.org/wiki/The_Making_of_the_Atomic_Bomb#Reception).
- Softwar: An Intimate Portrait of Larry Ellison and Oracle (2003) by Matthew Symonds. Ellison is founder of Oracle and the third-wealthiest person in the world.
  An interesting read! The book is very detailed. It feels like it covers every happening in Oracle since the founding until the publishing of the book.
  Thus I didn't end up finishing it, only about 250 out of 720 pages in.
  What's fun about the book is that Ellison required being able to add his own footnotes to it. Symonds had the last word on the text in the book itself, Ellison on the footnoes.
  These - surprisingly to me - made him come across as both humble and sympathetic.
- Rust in Action (2021) by Tim McNamara. Writing a lot of Rust at work, which is both fun and frustrating. Haven't used this enough to say anything sensible, though!
- Writing for Developers: Blogs that Get Read by Piotr Sarna, Cynthia Dunlop (2024). Very practical book about writing blogs. Was released late December, so haven't finished it yet. If you want to read it, you should join [Phil Eaton's book club on it](https://eatonphil.com/2025-writing-for-developers.html).

#### Source code for the graphic above

```
~ # cat book_stats.csv | awk -F ',' 'NR>1 {print $3}' | spark
▁▃█▅▃▄▃▁
```

where `spark` is Zach Holman's https://zachholman.com/spark/.
