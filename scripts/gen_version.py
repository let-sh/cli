import git
import os
repo = git.Repo(os.getcwd())
tags = sorted(repo.tags, key=lambda t: t.commit.committed_datetime)
tags = [str(tag) for tag in tags]
latest_tag = tags[-1][1:]
beta_tag = [tag for tag in list(tags) if ("beta" in tag)][-1][1:]
rc_tag = [tag for tag in list(tags) if ("rc" in tag)][-1][1:]
with open('version', 'w') as f:
    f.write(f'latest:{latest_tag}\nbeta:{beta_tag}\nrc:{rc_tag}')