ALTER TABLE ONLY public.notifications DROP CONSTRAINT notifications_user_id_fkey;
ALTER TABLE ONLY public.notifications DROP CONSTRAINT notifications_pkey;

DROP TABLE public.notifications;