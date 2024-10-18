CREATE DATABASE criticality_score;

\c criticality_score

-- ----------------------------
-- Table structure for arch_packages
-- ----------------------------
CREATE TABLE IF NOT EXISTS "public"."arch_packages" (
  "package" text COLLATE "pg_catalog"."default" NOT NULL,
  "version" text COLLATE "pg_catalog"."default",
  "depends_count" int8,
  "homepage" text COLLATE "pg_catalog"."default",
  "description" text COLLATE "pg_catalog"."default",
  "git_link" text COLLATE "pg_catalog"."default",
  "comment" text COLLATE "pg_catalog"."default",
  "link_confidence" float4
)
;

-- ----------------------------
-- Table structure for debian_packages
-- ----------------------------
CREATE TABLE IF NOT EXISTS "public"."debian_packages" (
  "package" text COLLATE "pg_catalog"."default" NOT NULL,
  "version" text COLLATE "pg_catalog"."default",
  "homepage" text COLLATE "pg_catalog"."default",
  "description" text COLLATE "pg_catalog"."default",
  "depends_count" int8,
  "git_link" text COLLATE "pg_catalog"."default",
  "comment" text COLLATE "pg_catalog"."default",
  "link_confidence" float4
)
;

-- ----------------------------
-- Table structure for git_metrics
-- ----------------------------
CREATE TABLE IF NOT EXISTS "public"."git_metrics" (
  "git_link" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "ecosystem" varchar(255) COLLATE "pg_catalog"."default",
  "from" int4 NOT NULL, -- from 0: package manager 1: enumerate_github
  "created_since" date,
  "updated_since" date,
  "contributor_count" int4,
  "commit_frequency" float8,
  "depsdev_count" int4,
  "deps_distro" float8,
  "scores" float8,
  "org_count" int4,
  "_name" varchar(255) COLLATE "pg_catalog"."default",
  "_owner" varchar(255) COLLATE "pg_catalog"."default",
  "_source" varchar(255) COLLATE "pg_catalog"."default"
)
;

-- ----------------------------
-- Table structure for nix_packages
-- ----------------------------
CREATE TABLE IF NOT EXISTS "public"."nix_packages" (
  "package" text COLLATE "pg_catalog"."default" NOT NULL,
  "version" text COLLATE "pg_catalog"."default",
  "depends_count" int8,
  "homepage" text COLLATE "pg_catalog"."default",
  "description" text COLLATE "pg_catalog"."default",
  "git_link" text COLLATE "pg_catalog"."default",
  "alias_link" text COLLATE "pg_catalog"."default",
  "link_confidence" float4
)
;

-- ----------------------------
-- Primary Key structure for table arch_packages
-- ----------------------------
ALTER TABLE "public"."arch_packages" ADD CONSTRAINT "arch_packages_pkey" PRIMARY KEY ("package");

-- ----------------------------
-- Primary Key structure for table debian_packages
-- ----------------------------
ALTER TABLE "public"."debian_packages" ADD CONSTRAINT "debian_packages_pkey" PRIMARY KEY ("package");

-- ----------------------------
-- Primary Key structure for table git_metrics
-- ----------------------------
ALTER TABLE "public"."git_metrics" ADD CONSTRAINT "git_metrics_pkey" PRIMARY KEY ("git_link");

-- ----------------------------
-- Primary Key structure for table git_metrics_backup
-- ----------------------------
ALTER TABLE "public"."git_metrics_backup" ADD CONSTRAINT "git_metrics_copy1_pkey" PRIMARY KEY ("git_link");

-- ----------------------------
-- Primary Key structure for table nix_packages
-- ----------------------------
ALTER TABLE "public"."nix_packages" ADD CONSTRAINT "arch_packages_copy1_pkey" PRIMARY KEY ("package");

