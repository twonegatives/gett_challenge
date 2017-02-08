--
-- PostgreSQL database dump
--

-- Dumped from database version 9.4.10
-- Dumped by pg_dump version 9.4.10
-- Started on 2017-02-10 01:04:07 MSK

SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;

--
-- TOC entry 1 (class 3079 OID 11897)
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- TOC entry 2052 (class 0 OID 0)
-- Dependencies: 1
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET search_path = public, pg_catalog;

--
-- TOC entry 176 (class 1259 OID 159624)
-- Name: driver_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE driver_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE driver_id_seq OWNER TO postgres;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- TOC entry 173 (class 1259 OID 159570)
-- Name: drivers; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE drivers (
    id integer DEFAULT nextval('driver_id_seq'::regclass) NOT NULL,
    name character varying(255) NOT NULL,
    license_number character varying(50) NOT NULL
);


ALTER TABLE drivers OWNER TO postgres;

--
-- TOC entry 175 (class 1259 OID 159613)
-- Name: metrics; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE metrics (
    id integer NOT NULL,
    metric_name character varying(255) NOT NULL,
    value integer NOT NULL,
    lat double precision NOT NULL,
    lon double precision NOT NULL,
    "timestamp" bigint NOT NULL,
    driver_id integer NOT NULL
);


ALTER TABLE metrics OWNER TO postgres;

--
-- TOC entry 2053 (class 0 OID 0)
-- Dependencies: 175
-- Name: TABLE metrics; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE metrics IS 'Thirteen decimal places will pin down the location to 111,111/10^13 = about 1 angstrom, around half the thickness of a small atom.
What do we need 18 digits contained in input json for?
http://gis.stackexchange.com/questions/8650/measuring-accuracy-of-latitude-and-longitude';


--
-- TOC entry 174 (class 1259 OID 159611)
-- Name: metrics_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE metrics_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE metrics_id_seq OWNER TO postgres;

--
-- TOC entry 2054 (class 0 OID 0)
-- Dependencies: 174
-- Name: metrics_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE metrics_id_seq OWNED BY metrics.id;


--
-- TOC entry 1928 (class 2604 OID 159616)
-- Name: id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY metrics ALTER COLUMN id SET DEFAULT nextval('metrics_id_seq'::regclass);


--
-- TOC entry 1931 (class 2606 OID 159574)
-- Name: drivers_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY drivers
    ADD CONSTRAINT drivers_pkey PRIMARY KEY (id);


--
-- TOC entry 1935 (class 2606 OID 159618)
-- Name: metrics_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY metrics
    ADD CONSTRAINT metrics_pkey PRIMARY KEY (id);


--
-- TOC entry 1929 (class 1259 OID 159665)
-- Name: drivers_id_idx; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE INDEX drivers_id_idx ON drivers USING btree (id);


--
-- TOC entry 1932 (class 1259 OID 159667)
-- Name: metrics_driver_id_idx; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE INDEX metrics_driver_id_idx ON metrics USING btree (driver_id);


--
-- TOC entry 1933 (class 1259 OID 159666)
-- Name: metrics_id_idx; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE INDEX metrics_id_idx ON metrics USING btree (id);


--
-- TOC entry 2051 (class 0 OID 0)
-- Dependencies: 6
-- Name: public; Type: ACL; Schema: -; Owner: postgres
--

REVOKE ALL ON SCHEMA public FROM PUBLIC;
REVOKE ALL ON SCHEMA public FROM postgres;
GRANT ALL ON SCHEMA public TO postgres;
GRANT ALL ON SCHEMA public TO PUBLIC;


-- Completed on 2017-02-10 01:04:07 MSK

--
-- PostgreSQL database dump complete
--

