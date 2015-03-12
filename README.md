# RightScale events RSS watcher

This simple tool polls the RightScale RSS feed for events on your account and
displays them sorted by time, showing you only the latest set. You could compare
it to a `tail -f` for Rightscale.

Whenever I'm working on a process which depends on some asyncronous process
within RightScale to complete, I like to run this tool in it's own window and
just let it run all day.

Install it using `go get`:

```shell
go get -u -v github.com/rtlong/rightscale-rss-events-log
```

Then run it, passing it your RightScale RSS URL as the only argument:

```shell
rightscale-rss-events-log ${URL}
```

See more about where to get the URL here:
[RightScale 101 > Management Tools > Events](https://support.rightscale.com/12-Guides/RightScale_101/08-Management_Tools/Events)

The URL should look similar to this:

```shell
https://us-4.rightscale.com/acct/${account_id}/user_notifications/feed.atom?feed_token=${token}
```
