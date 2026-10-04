CREATE TABLE public.users (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    email character varying(255) NOT NULL,
    password character varying NOT NULL,
    name character varying NOT NULL,
    avatar_url character varying(255) DEFAULT NULL::character varying,
    location character varying(255) DEFAULT NULL::character varying,
    bio character varying(255) DEFAULT NULL::character varying,
    job character varying(255) DEFAULT NULL::character varying,
    role character varying DEFAULT 'attendee'::character varying NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);