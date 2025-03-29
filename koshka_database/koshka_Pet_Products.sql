--
-- PostgreSQL database dump
--

-- Dumped from database version 17.4 (Debian 17.4-1.pgdg120+2)
-- Dumped by pg_dump version 17.4 (Debian 17.4-1.pgdg120+2)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: api_keys; Type: TABLE; Schema: public; Owner: KOSHKA
--

CREATE TABLE public.api_keys (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    created_by character varying(255),
    app character varying(50) NOT NULL,
    api_key text DEFAULT encode(public.gen_random_bytes(32), 'hex'::text) NOT NULL,
    created_at timestamp without time zone DEFAULT now()
);


ALTER TABLE public.api_keys OWNER TO "KOSHKA";

--
-- Name: api_permissions; Type: TABLE; Schema: public; Owner: KOSHKA
--

CREATE TABLE public.api_permissions (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    api_key_id uuid,
    route_pattern text NOT NULL,
    method text NOT NULL,
    allowed boolean DEFAULT false,
    CONSTRAINT api_permissions_method_check CHECK ((method = ANY (ARRAY['GET'::text, 'POST'::text, 'PUT'::text, 'DELETE'::text])))
);


ALTER TABLE public.api_permissions OWNER TO "KOSHKA";

--
-- Name: images; Type: TABLE; Schema: public; Owner: KOSHKA
--

CREATE TABLE public.images (
    id integer NOT NULL,
    hash text NOT NULL,
    file text NOT NULL
);


ALTER TABLE public.images OWNER TO "KOSHKA";

--
-- Name: images_id_seq; Type: SEQUENCE; Schema: public; Owner: KOSHKA
--

CREATE SEQUENCE public.images_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.images_id_seq OWNER TO "KOSHKA";

--
-- Name: images_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: KOSHKA
--

ALTER SEQUENCE public.images_id_seq OWNED BY public.images.id;


--
-- Name: personas; Type: TABLE; Schema: public; Owner: KOSHKA
--

CREATE TABLE public.personas (
    id integer NOT NULL,
    name character varying(10) NOT NULL,
    description character varying(50),
    admin boolean DEFAULT false
);


ALTER TABLE public.personas OWNER TO "KOSHKA";

--
-- Name: personas_id_seq; Type: SEQUENCE; Schema: public; Owner: KOSHKA
--

CREATE SEQUENCE public.personas_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.personas_id_seq OWNER TO "KOSHKA";

--
-- Name: personas_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: KOSHKA
--

ALTER SEQUENCE public.personas_id_seq OWNED BY public.personas.id;


--
-- Name: pet_collars; Type: TABLE; Schema: public; Owner: KOSHKA
--

CREATE TABLE public.pet_collars (
    pet_id uuid,
    collar_id character varying(6) NOT NULL,
    assigned_at timestamp without time zone DEFAULT now(),
    user_id character varying NOT NULL
);


ALTER TABLE public.pet_collars OWNER TO "KOSHKA";

--
-- Name: pet_profiles; Type: TABLE; Schema: public; Owner: KOSHKA
--

CREATE TABLE public.pet_profiles (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    picture_id integer NOT NULL,
    name character varying(15) NOT NULL,
    date_of_birth date,
    tag_id character varying(6),
    city_licence character varying(10),
    neutered boolean DEFAULT false,
    vaccinated boolean DEFAULT false,
    illnesses_code character varying(6),
    updated_at timestamp without time zone DEFAULT now()
);


ALTER TABLE public.pet_profiles OWNER TO "KOSHKA";

--
-- Name: product_features; Type: TABLE; Schema: public; Owner: KOSHKA
--

CREATE TABLE public.product_features (
    id integer NOT NULL,
    text character varying(30) NOT NULL
);


ALTER TABLE public.product_features OWNER TO "KOSHKA";

--
-- Name: product_features_id_seq; Type: SEQUENCE; Schema: public; Owner: KOSHKA
--

CREATE SEQUENCE public.product_features_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.product_features_id_seq OWNER TO "KOSHKA";

--
-- Name: product_features_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: KOSHKA
--

ALTER SEQUENCE public.product_features_id_seq OWNED BY public.product_features.id;


--
-- Name: product_features_map; Type: TABLE; Schema: public; Owner: KOSHKA
--

CREATE TABLE public.product_features_map (
    product_sku character varying(16) NOT NULL,
    feature_id integer NOT NULL
);


ALTER TABLE public.product_features_map OWNER TO "KOSHKA";

--
-- Name: products; Type: TABLE; Schema: public; Owner: KOSHKA
--

CREATE TABLE public.products (
    sku character varying(16) NOT NULL,
    description text,
    image_id integer
);


ALTER TABLE public.products OWNER TO "KOSHKA";

--
-- Name: purchase_history; Type: TABLE; Schema: public; Owner: KOSHKA
--

CREATE TABLE public.purchase_history (
    product_sku character varying(16) NOT NULL,
    registration_date date NOT NULL,
    tag_id character varying(6) NOT NULL
);


ALTER TABLE public.purchase_history OWNER TO "KOSHKA";

--
-- Name: reunite_collars; Type: TABLE; Schema: public; Owner: KOSHKA
--

CREATE TABLE public.reunite_collars (
    tag_id character varying(6) NOT NULL,
    product_sku character varying(16),
    active boolean DEFAULT true,
    registered boolean DEFAULT false
);


ALTER TABLE public.reunite_collars OWNER TO "KOSHKA";

--
-- Name: user_pets; Type: TABLE; Schema: public; Owner: KOSHKA
--

CREATE TABLE public.user_pets (
    user_id character varying(50) NOT NULL,
    pet_id uuid NOT NULL
);


ALTER TABLE public.user_pets OWNER TO "KOSHKA";

--
-- Name: user_profiles; Type: TABLE; Schema: public; Owner: KOSHKA
--

CREATE TABLE public.user_profiles (
    email character varying(50) NOT NULL,
    name character varying(50) NOT NULL,
    address text NOT NULL,
    mobile_number character varying(15) NOT NULL,
    purchase_history integer[] DEFAULT '{}'::integer[],
    reunite_collars integer[] DEFAULT '{}'::integer[],
    pets integer[] DEFAULT '{}'::integer[],
    updated_at timestamp without time zone DEFAULT now()
);


ALTER TABLE public.user_profiles OWNER TO "KOSHKA";

--
-- Name: users; Type: TABLE; Schema: public; Owner: KOSHKA
--

CREATE TABLE public.users (
    email character varying(50) NOT NULL,
    password text NOT NULL,
    salt text NOT NULL,
    persona_id integer DEFAULT 2,
    marketing boolean DEFAULT true,
    cookies json,
    verified boolean DEFAULT false,
    locked boolean DEFAULT false,
    verification_token text,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);


ALTER TABLE public.users OWNER TO "KOSHKA";

--
-- Name: veterinary_illnesses; Type: TABLE; Schema: public; Owner: KOSHKA
--

CREATE TABLE public.veterinary_illnesses (
    code character varying(6) NOT NULL,
    name character varying(30) NOT NULL,
    description text
);


ALTER TABLE public.veterinary_illnesses OWNER TO "KOSHKA";

--
-- Name: images id; Type: DEFAULT; Schema: public; Owner: KOSHKA
--

ALTER TABLE ONLY public.images ALTER COLUMN id SET DEFAULT nextval('public.images_id_seq'::regclass);


--
-- Name: personas id; Type: DEFAULT; Schema: public; Owner: KOSHKA
--

ALTER TABLE ONLY public.personas ALTER COLUMN id SET DEFAULT nextval('public.personas_id_seq'::regclass);


--
-- Name: product_features id; Type: DEFAULT; Schema: public; Owner: KOSHKA
--

ALTER TABLE ONLY public.product_features ALTER COLUMN id SET DEFAULT nextval('public.product_features_id_seq'::regclass);


--
-- Name: api_keys api_keys_api_key_key; Type: CONSTRAINT; Schema: public; Owner: KOSHKA
--

ALTER TABLE ONLY public.api_keys
    ADD CONSTRAINT api_keys_api_key_key UNIQUE (api_key);


--
-- Name: api_keys api_keys_pkey; Type: CONSTRAINT; Schema: public; Owner: KOSHKA
--

ALTER TABLE ONLY public.api_keys
    ADD CONSTRAINT api_keys_pkey PRIMARY KEY (id);


--
-- Name: api_permissions api_permissions_api_key_id_route_pattern_method_key; Type: CONSTRAINT; Schema: public; Owner: KOSHKA
--

ALTER TABLE ONLY public.api_permissions
    ADD CONSTRAINT api_permissions_api_key_id_route_pattern_method_key UNIQUE (api_key_id, route_pattern, method);


--
-- Name: api_permissions api_permissions_pkey; Type: CONSTRAINT; Schema: public; Owner: KOSHKA
--

ALTER TABLE ONLY public.api_permissions
    ADD CONSTRAINT api_permissions_pkey PRIMARY KEY (id);


--
-- Name: images images_hash_key; Type: CONSTRAINT; Schema: public; Owner: KOSHKA
--

ALTER TABLE ONLY public.images
    ADD CONSTRAINT images_hash_key UNIQUE (hash);


--
-- Name: images images_pkey; Type: CONSTRAINT; Schema: public; Owner: KOSHKA
--

ALTER TABLE ONLY public.images
    ADD CONSTRAINT images_pkey PRIMARY KEY (id);


--
-- Name: personas personas_name_key; Type: CONSTRAINT; Schema: public; Owner: KOSHKA
--

ALTER TABLE ONLY public.personas
    ADD CONSTRAINT personas_name_key UNIQUE (name);


--
-- Name: personas personas_pkey; Type: CONSTRAINT; Schema: public; Owner: KOSHKA
--

ALTER TABLE ONLY public.personas
    ADD CONSTRAINT personas_pkey PRIMARY KEY (id);


--
-- Name: pet_collars pet_collars_unique; Type: CONSTRAINT; Schema: public; Owner: KOSHKA
--

ALTER TABLE ONLY public.pet_collars
    ADD CONSTRAINT pet_collars_unique UNIQUE (collar_id, user_id);


--
-- Name: pet_profiles pet_profiles_pkey; Type: CONSTRAINT; Schema: public; Owner: KOSHKA
--

ALTER TABLE ONLY public.pet_profiles
    ADD CONSTRAINT pet_profiles_pkey PRIMARY KEY (id);


--
-- Name: product_features_map product_features_map_pkey; Type: CONSTRAINT; Schema: public; Owner: KOSHKA
--

ALTER TABLE ONLY public.product_features_map
    ADD CONSTRAINT product_features_map_pkey PRIMARY KEY (product_sku, feature_id);


--
-- Name: product_features product_features_pkey; Type: CONSTRAINT; Schema: public; Owner: KOSHKA
--

ALTER TABLE ONLY public.product_features
    ADD CONSTRAINT product_features_pkey PRIMARY KEY (id);


--
-- Name: products products_pkey; Type: CONSTRAINT; Schema: public; Owner: KOSHKA
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_pkey PRIMARY KEY (sku);


--
-- Name: purchase_history purchase_history_pkey; Type: CONSTRAINT; Schema: public; Owner: KOSHKA
--

ALTER TABLE ONLY public.purchase_history
    ADD CONSTRAINT purchase_history_pkey PRIMARY KEY (product_sku, tag_id);


--
-- Name: reunite_collars reunite_collars_pkey; Type: CONSTRAINT; Schema: public; Owner: KOSHKA
--

ALTER TABLE ONLY public.reunite_collars
    ADD CONSTRAINT reunite_collars_pkey PRIMARY KEY (tag_id);


--
-- Name: api_permissions unique_route_method; Type: CONSTRAINT; Schema: public; Owner: KOSHKA
--

ALTER TABLE ONLY public.api_permissions
    ADD CONSTRAINT unique_route_method UNIQUE (route_pattern, method);


--
-- Name: user_pets user_pets_pkey; Type: CONSTRAINT; Schema: public; Owner: KOSHKA
--

ALTER TABLE ONLY public.user_pets
    ADD CONSTRAINT user_pets_pkey PRIMARY KEY (user_id, pet_id);


--
-- Name: user_profiles user_profiles_pkey; Type: CONSTRAINT; Schema: public; Owner: KOSHKA
--

ALTER TABLE ONLY public.user_profiles
    ADD CONSTRAINT user_profiles_pkey PRIMARY KEY (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: KOSHKA
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (email);


--
-- Name: veterinary_illnesses veterinary_illnesses_pkey; Type: CONSTRAINT; Schema: public; Owner: KOSHKA
--

ALTER TABLE ONLY public.veterinary_illnesses
    ADD CONSTRAINT veterinary_illnesses_pkey PRIMARY KEY (code);


--
-- Name: idx_images_hash; Type: INDEX; Schema: public; Owner: KOSHKA
--

CREATE INDEX idx_images_hash ON public.images USING btree (hash);


--
-- Name: pet_profiles trigger_pet_profiles_updated_at; Type: TRIGGER; Schema: public; Owner: KOSHKA
--

CREATE TRIGGER trigger_pet_profiles_updated_at BEFORE UPDATE ON public.pet_profiles FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: user_profiles trigger_user_profiles_updated_at; Type: TRIGGER; Schema: public; Owner: KOSHKA
--

CREATE TRIGGER trigger_user_profiles_updated_at BEFORE UPDATE ON public.user_profiles FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: api_keys api_keys_created_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: KOSHKA
--

ALTER TABLE ONLY public.api_keys
    ADD CONSTRAINT api_keys_created_by_fkey FOREIGN KEY (created_by) REFERENCES public.users(email) ON DELETE CASCADE;


--
-- Name: api_permissions api_permissions_api_key_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: KOSHKA
--

ALTER TABLE ONLY public.api_permissions
    ADD CONSTRAINT api_permissions_api_key_id_fkey FOREIGN KEY (api_key_id) REFERENCES public.api_keys(id) ON DELETE CASCADE;


--
-- Name: pet_collars pet_collars_collar_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: KOSHKA
--

ALTER TABLE ONLY public.pet_collars
    ADD CONSTRAINT pet_collars_collar_id_fkey FOREIGN KEY (collar_id) REFERENCES public.reunite_collars(tag_id) ON DELETE SET NULL;


--
-- Name: pet_collars pet_collars_pet_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: KOSHKA
--

ALTER TABLE ONLY public.pet_collars
    ADD CONSTRAINT pet_collars_pet_id_fkey FOREIGN KEY (pet_id) REFERENCES public.pet_profiles(id) ON DELETE SET NULL;


--
-- Name: pet_collars pet_collars_user_profiles_fk; Type: FK CONSTRAINT; Schema: public; Owner: KOSHKA
--

ALTER TABLE ONLY public.pet_collars
    ADD CONSTRAINT pet_collars_user_profiles_fk FOREIGN KEY (user_id) REFERENCES public.user_profiles(email) ON DELETE CASCADE;


--
-- Name: pet_profiles pet_profiles_illnesses_code_fkey; Type: FK CONSTRAINT; Schema: public; Owner: KOSHKA
--

ALTER TABLE ONLY public.pet_profiles
    ADD CONSTRAINT pet_profiles_illnesses_code_fkey FOREIGN KEY (illnesses_code) REFERENCES public.veterinary_illnesses(code);


--
-- Name: pet_profiles pet_profiles_picture_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: KOSHKA
--

ALTER TABLE ONLY public.pet_profiles
    ADD CONSTRAINT pet_profiles_picture_id_fkey FOREIGN KEY (picture_id) REFERENCES public.images(id);


--
-- Name: pet_profiles pet_profiles_tag_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: KOSHKA
--

ALTER TABLE ONLY public.pet_profiles
    ADD CONSTRAINT pet_profiles_tag_id_fkey FOREIGN KEY (tag_id) REFERENCES public.reunite_collars(tag_id);


--
-- Name: product_features_map product_features_map_feature_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: KOSHKA
--

ALTER TABLE ONLY public.product_features_map
    ADD CONSTRAINT product_features_map_feature_id_fkey FOREIGN KEY (feature_id) REFERENCES public.product_features(id);


--
-- Name: product_features_map product_features_map_product_sku_fkey; Type: FK CONSTRAINT; Schema: public; Owner: KOSHKA
--

ALTER TABLE ONLY public.product_features_map
    ADD CONSTRAINT product_features_map_product_sku_fkey FOREIGN KEY (product_sku) REFERENCES public.products(sku);


--
-- Name: products products_image_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: KOSHKA
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_image_id_fkey FOREIGN KEY (image_id) REFERENCES public.images(id);


--
-- Name: purchase_history purchase_history_product_sku_fkey; Type: FK CONSTRAINT; Schema: public; Owner: KOSHKA
--

ALTER TABLE ONLY public.purchase_history
    ADD CONSTRAINT purchase_history_product_sku_fkey FOREIGN KEY (product_sku) REFERENCES public.products(sku);


--
-- Name: purchase_history purchase_history_tag_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: KOSHKA
--

ALTER TABLE ONLY public.purchase_history
    ADD CONSTRAINT purchase_history_tag_id_fkey FOREIGN KEY (tag_id) REFERENCES public.reunite_collars(tag_id);


--
-- Name: reunite_collars reunite_collars_product_sku_fkey; Type: FK CONSTRAINT; Schema: public; Owner: KOSHKA
--

ALTER TABLE ONLY public.reunite_collars
    ADD CONSTRAINT reunite_collars_product_sku_fkey FOREIGN KEY (product_sku) REFERENCES public.products(sku);


--
-- Name: user_pets user_pets_pet_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: KOSHKA
--

ALTER TABLE ONLY public.user_pets
    ADD CONSTRAINT user_pets_pet_id_fkey FOREIGN KEY (pet_id) REFERENCES public.pet_profiles(id);


--
-- Name: user_pets user_pets_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: KOSHKA
--

ALTER TABLE ONLY public.user_pets
    ADD CONSTRAINT user_pets_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(email) ON DELETE CASCADE;


--
-- Name: user_profiles user_profiles_email_fkey; Type: FK CONSTRAINT; Schema: public; Owner: KOSHKA
--

ALTER TABLE ONLY public.user_profiles
    ADD CONSTRAINT user_profiles_email_fkey FOREIGN KEY (email) REFERENCES public.users(email);


--
-- Name: users users_persona_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: KOSHKA
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_persona_id_fkey FOREIGN KEY (persona_id) REFERENCES public.personas(id);


--
-- PostgreSQL database dump complete
--

