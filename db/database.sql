/*
 Navicat Premium Data Transfer

 Source Server         : kesehatan
 Source Server Type    : PostgreSQL
 Source Server Version : 140005 (140005)
 Source Host           : localhost:5432
 Source Catalog        : perpustakaan-golang
 Source Schema         : public

 Target Server Type    : PostgreSQL
 Target Server Version : 140005 (140005)
 File Encoding         : 65001

 Date: 10/02/2025 20:52:20
*/


-- ----------------------------
-- Sequence structure for user_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."user_id_seq";
CREATE SEQUENCE "public"."user_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Table structure for master_buku
-- ----------------------------
DROP TABLE IF EXISTS "public"."master_buku";
CREATE TABLE "public"."master_buku" (
  "id" uuid NOT NULL,
  "judul" varchar(255) COLLATE "pg_catalog"."default",
  "pengarang" varchar(255) COLLATE "pg_catalog"."default",
  "penerbit" varchar(255) COLLATE "pg_catalog"."default",
  "isbn" varchar(17) COLLATE "pg_catalog"."default",
  "tahun_terbit" int8,
  "kategori" varchar(255) COLLATE "pg_catalog"."default",
  "deskripsi" text COLLATE "pg_catalog"."default",
  "foto" varchar(255) COLLATE "pg_catalog"."default",
  "jumlah_eksemplar" int8,
  "jumlah_ketersediaan_eksemplar" int8,
  "created_at" date,
  "updated_at" date
)
;

-- ----------------------------
-- Table structure for transaksi_buku
-- ----------------------------
DROP TABLE IF EXISTS "public"."transaksi_buku";
CREATE TABLE "public"."transaksi_buku" (
  "id" uuid NOT NULL,
  "id_anggota" int4,
  "id_buku" uuid,
  "tanggal_pinjam" timestamp(6),
  "tanggal_jatuh_tempo" timestamp(6),
  "tanggal_kembali" timestamp(6),
  "denda" numeric(255,0),
  "status" int2
)
;

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS "public"."user";
CREATE TABLE "public"."user" (
  "id" int4 NOT NULL DEFAULT nextval('user_id_seq'::regclass),
  "email" varchar COLLATE "pg_catalog"."default",
  "password" varchar COLLATE "pg_catalog"."default",
  "nama" varchar COLLATE "pg_catalog"."default",
  "updated_at" int4,
  "created_at" int4,
  "alamat" text COLLATE "pg_catalog"."default",
  "nomor_telepon" varchar(20) COLLATE "pg_catalog"."default",
  "tanggal_lahir" date,
  "tanggal_join" date,
  "status_anggota" int2,
  "foto" varchar COLLATE "pg_catalog"."default"
)
;

-- ----------------------------
-- Function structure for created_at_column
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."created_at_column"();
CREATE OR REPLACE FUNCTION "public"."created_at_column"()
  RETURNS "pg_catalog"."trigger" AS $BODY$

BEGIN
	NEW.updated_at = EXTRACT(EPOCH FROM NOW());
	NEW.created_at = EXTRACT(EPOCH FROM NOW());
    RETURN NEW;
END;

$BODY$
  LANGUAGE plpgsql VOLATILE
  COST 100;

-- ----------------------------
-- Function structure for update_at_column
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."update_at_column"();
CREATE OR REPLACE FUNCTION "public"."update_at_column"()
  RETURNS "pg_catalog"."trigger" AS $BODY$

BEGIN
    NEW.updated_at = EXTRACT(EPOCH FROM NOW());
    RETURN NEW;
END;

$BODY$
  LANGUAGE plpgsql VOLATILE
  COST 100;

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."user_id_seq"
OWNED BY "public"."user"."id";
SELECT setval('"public"."user_id_seq"', 1, false);

-- ----------------------------
-- Primary Key structure for table master_buku
-- ----------------------------
ALTER TABLE "public"."master_buku" ADD CONSTRAINT "master_buku_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table transaksi_buku
-- ----------------------------
ALTER TABLE "public"."transaksi_buku" ADD CONSTRAINT "transaksi_buku_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Triggers structure for table user
-- ----------------------------
CREATE TRIGGER "create_user_created_at" BEFORE INSERT ON "public"."user"
FOR EACH ROW
EXECUTE PROCEDURE "public"."created_at_column"();
CREATE TRIGGER "update_user_updated_at" BEFORE UPDATE ON "public"."user"
FOR EACH ROW
EXECUTE PROCEDURE "public"."update_at_column"();

-- ----------------------------
-- Primary Key structure for table user
-- ----------------------------
ALTER TABLE "public"."user" ADD CONSTRAINT "user_id" PRIMARY KEY ("id");

-- ----------------------------
-- Foreign Keys structure for table transaksi_buku
-- ----------------------------
ALTER TABLE "public"."transaksi_buku" ADD CONSTRAINT "foreign_anggota" FOREIGN KEY ("id_anggota") REFERENCES "public"."user" ("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "public"."transaksi_buku" ADD CONSTRAINT "foreign_buku" FOREIGN KEY ("id_buku") REFERENCES "public"."master_buku" ("id") ON DELETE NO ACTION ON UPDATE CASCADE;
