CREATE DATABASE criticality_score;

\c criticality_score

create table if not exists public.arch_packages
(
    package       text not null
        primary key,
    version       text,
    depends_count bigint,
    homepage      text,
    description   text,
    git_link      text,
    comment       text
);

create table if not exists public.debian_packages
(
    package       text not null
        primary key,
    version       text,
    homepage      text,
    description   text,
    depends_count bigint,
    git_link      text,
    comment       text
);

create table if not exists public.git_metrics
(
    git_link          varchar(255) not null
        primary key,
    ecosystem         varchar(255),
    "from"            integer      not null,
    created_since     date,
    updated_since     date,
    contributor_count integer          default 0,
    commit_frequency  double precision default 0,
    depsdev_count     integer          default 0,
    dist_impact       double precision default 0,
    scores            double precision default 0,
    org_count         integer          default 0,
    _name             varchar(255),
    _owner            varchar(255),
    _source           varchar(255),
    need_update       boolean,
    license           varchar(255)
);

create table if not exists public.nix_packages
(
    package       text not null
        constraint arch_packages_copy1_pkey
            primary key,
    version       text,
    depends_count bigint,
    homepage      text,
    description   text,
    git_link      text,
    alias_link    text
);

create table if not exists public.homebrew_packages
(
    package       text not null
        constraint debian_packages_copy1_pkey
            primary key,
    homepage      text,
    description   text,
    depends_count bigint,
    git_link      text,
    alias_link    text
);

create table if not exists public.github_links
(
    git_link text not null
        primary key
);

create table if not exists public.arch_relationships
(
    frompackage varchar(255) not null
        references public.arch_packages,
    topackage   varchar(255) not null,
    primary key (frompackage, topackage)
);

create table if not exists public.debian_relationships
(
    frompackage varchar(255) not null
        references public.debian_packages,
    topackage   varchar(255) not null,
    primary key (frompackage, topackage)
);

create table if not exists public.gentoo_packages
(
    package         text not null
        constraint homebrew_packages_copy1_pkey
            primary key,
    version         text,
    homepage        text,
    description     text,
    depends_count   bigint,
    git_link        text,
    alias_link      text,
    link_confidence real
);

create table if not exists public.gentoo_relationships
(
    frompackage varchar(255) not null
        references public.gentoo_packages,
    topackage   varchar(255) not null,
    primary key (frompackage, topackage)
);

create table if not exists public.homebrew_relationships
(
    frompackage varchar(255) not null
        references public.homebrew_packages,
    topackage   varchar(255) not null,
    primary key (frompackage, topackage)
);

create table if not exists public.nix_relationships
(
    frompackage varchar(255) not null
        references public.nix_packages,
    topackage   varchar(255) not null,
    primary key (frompackage, topackage)
);

create table if not exists public.git_repositories
(
    git_link varchar(255) not null
        primary key,
    industry integer,
    domestic boolean
);

