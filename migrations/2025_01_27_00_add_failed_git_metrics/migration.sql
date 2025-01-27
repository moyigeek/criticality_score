create table failed_git_metrics (
    git_link text primary key,
    message text,
    updated_time timestamp,
    times int
);
