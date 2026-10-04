CREATE TABLE public.tags (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name character varying NOT NULL
);

ALTER TABLE ONLY public.tags
    ADD CONSTRAINT tags_pkey PRIMARY KEY (id);