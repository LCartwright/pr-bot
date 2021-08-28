# PR-bot
[P]ull [R]equest bot - produces a slack formatted summary (to stdout) on open pull requests, with a mapping from GitHub login to slack mention username

## Sample output

> *Pull Requests for [TOP](https://github.com/Paymentsense/TOP/pulls?q=is%3Apr+is%3Aopen+-label%3Awip+review%3Arequired+sort%3Acreated-asc)*
> 
> _Beep, boop, Happy Sunday! We have *6* Open Pull Requests today!_
> * *[[AG-557] Metrics package with Instrumentation for HTTP (mux) routing](https://github.com/Paymentsense/TOP/pull/707)*
>     * Author: @laurence , open for: *62d0h*, last updated: *Aug 27 17:17:30*
>     * Requested reviewer(s): @Richard Lloyd , @fabio.ornellas .
> * *[[AG-245] HSM service v2 (Part 1)](https://github.com/Paymentsense/TOP/pull/832)*
>     * Author: @pedro , open for: *12d9h*, last updated: *Aug 27 14:08:07*
>     * Requested reviewer(s): @Ben K , @dale .
> * *[bugfix/ [AG-495] Missing Client.Close() throughout code base](https://github.com/Paymentsense/TOP/pull/833)*
>     * Author: @Kingsley Edore , open for: *12d4h*, last updated: *Aug 27 16:47:14*
>     * Requested reviewer(s): @alex.mckinlay , @Harry Tennent .
> * *[Feature/ag 548 create anti-DoS IP list for app-eng firewall](https://github.com/Paymentsense/TOP/pull/861)*
>     * Author: [evansecarch](https://github.com/evansecarch), open for: *5d0h*, last updated: *Aug 25 07:30:41*
>     * *No pending reviewers!*
> * *[[AG-580] enable Skaffold for tokenization-service](https://github.com/Paymentsense/TOP/pull/866)*
>     * Author: @davidgw , open for: *3d22h*, last updated: *Aug 27 17:07:34*
>     * Requested reviewer(s): @Ben K , @laurence , @fabio.ornellas .
> * *[[AG-553] Remove terminal-orders-status-updates-topic from Terraform Files](https://github.com/Paymentsense/TOP/pull/869)*
>     * Author: @Kingsley Edore , open for: *2d6h*, last updated: *Aug 27 15:52:31*
>    * Requested reviewer(s): @Richard Lloyd .

## Confirguration
Uses [dotenv](https://github.com/joho/godotenv) and [envconfig](https://github.com/kelseyhightower/envconfig) to bind the following environment variables
* GH_ACCESS_TOKEN - GitHub access token with *repo* read privs
* GH_OWNER - Owner e.g. *LCartwright*
* GH_REPO - Repo e.g. *pr-bot*
* GH_URL - URL e.g. *https://github.com/LCartwright/pr-bot*
* GH_BASE_BRANCH - e.g. *main*
* GH_LOGIN_TO_SLACK_USERNAMES - Map in the following format <GitHub Login>:<Slack Mention>,... e.g. *LCartwright:laurence,Octocat:Octo Cat*

## TODO
* Integrate with Slack
* Resolve PR status