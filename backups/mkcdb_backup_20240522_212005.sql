PGDMP                         |            mkcdb    14.1    15.6 (Debian 15.6-0+deb12u1) )    H           0    0    ENCODING    ENCODING        SET client_encoding = 'UTF8';
                      false            I           0    0 
   STDSTRINGS 
   STDSTRINGS     (   SET standard_conforming_strings = 'on';
                      false            J           0    0 
   SEARCHPATH 
   SEARCHPATH     8   SELECT pg_catalog.set_config('search_path', '', false);
                      false            K           1262    16384    mkcdb    DATABASE     p   CREATE DATABASE mkcdb WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE_PROVIDER = libc LOCALE = 'en_US.utf8';
    DROP DATABASE mkcdb;
                postgres    false                        2615    2200    public    SCHEMA     2   -- *not* creating schema, since initdb creates it
 2   -- *not* dropping schema, since initdb creates it
                postgres    false            L           0    0    SCHEMA public    ACL     Q   REVOKE USAGE ON SCHEMA public FROM PUBLIC;
GRANT ALL ON SCHEMA public TO PUBLIC;
                   postgres    false    5                        3079    16393 	   uuid-ossp 	   EXTENSION     ?   CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;
    DROP EXTENSION "uuid-ossp";
                   false    5            M           0    0    EXTENSION "uuid-ossp"    COMMENT     W   COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';
                        false    2            �            1259    46574    administrator    TABLE     ;  CREATE TABLE public.administrator (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name character varying(255) DEFAULT ''::character varying NOT NULL,
    email character varying(255) DEFAULT ''::character varying NOT NULL,
    password character varying(255) DEFAULT ''::character varying NOT NULL
);
 !   DROP TABLE public.administrator;
       public         heap    postgres    false    2    5    5            �            1259    46556    customer    TABLE     +  CREATE TABLE public.customer (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    first_name character varying(255) DEFAULT ''::character varying NOT NULL,
    second_name character varying(255) DEFAULT ''::character varying NOT NULL,
    login character varying(255) DEFAULT ''::character varying NOT NULL,
    password character varying(255) DEFAULT ''::character varying NOT NULL,
    email character varying(255) DEFAULT ''::character varying NOT NULL,
    type integer DEFAULT 0 NOT NULL,
    subscription_end timestamp without time zone
);
    DROP TABLE public.customer;
       public         heap    postgres    false    2    5    5            �            1259    46620    file    TABLE     �  CREATE TABLE public.file (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    project_id uuid,
    filename character varying(255) DEFAULT ''::character varying NOT NULL,
    extension character varying(255) DEFAULT ''::character varying NOT NULL,
    size integer DEFAULT 0 NOT NULL,
    file_path character varying(255) DEFAULT ''::character varying NOT NULL,
    update_datetime timestamp without time zone
);
    DROP TABLE public.file;
       public         heap    postgres    false    2    5    5            �            1259    16386    goose_db_version    TABLE     �   CREATE TABLE public.goose_db_version (
    id integer NOT NULL,
    version_id bigint NOT NULL,
    is_applied boolean NOT NULL,
    tstamp timestamp without time zone DEFAULT now()
);
 $   DROP TABLE public.goose_db_version;
       public         heap    postgres    false    5            �            1259    16385    goose_db_version_id_seq    SEQUENCE     �   CREATE SEQUENCE public.goose_db_version_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 .   DROP SEQUENCE public.goose_db_version_id_seq;
       public          postgres    false    211    5            N           0    0    goose_db_version_id_seq    SEQUENCE OWNED BY     S   ALTER SEQUENCE public.goose_db_version_id_seq OWNED BY public.goose_db_version.id;
          public          postgres    false    210            �            1259    46637    note    TABLE     i  CREATE TABLE public.note (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    project_id uuid,
    title character varying(255) DEFAULT ''::character varying NOT NULL,
    content character varying DEFAULT ''::character varying NOT NULL,
    update_datetime timestamp without time zone,
    deadline timestamp without time zone,
    overdue integer
);
    DROP TABLE public.note;
       public         heap    postgres    false    2    5    5            �            1259    46585    project    TABLE       CREATE TABLE public.project (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    owner_id uuid,
    capacity integer DEFAULT 0 NOT NULL,
    name character varying(255) DEFAULT ''::character varying NOT NULL,
    creation_date timestamp without time zone,
    admin_id uuid
);
    DROP TABLE public.project;
       public         heap    postgres    false    2    5    5            �            1259    46603    project_access    TABLE     �   CREATE TABLE public.project_access (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    project_id uuid,
    customer_id uuid,
    customer_access integer DEFAULT 0 NOT NULL
);
 "   DROP TABLE public.project_access;
       public         heap    postgres    false    2    5    5            �           2604    16389    goose_db_version id    DEFAULT     z   ALTER TABLE ONLY public.goose_db_version ALTER COLUMN id SET DEFAULT nextval('public.goose_db_version_id_seq'::regclass);
 B   ALTER TABLE public.goose_db_version ALTER COLUMN id DROP DEFAULT;
       public          postgres    false    211    210    211            A          0    46574    administrator 
   TABLE DATA           B   COPY public.administrator (id, name, email, password) FROM stdin;
    public          postgres    false    213   M4       @          0    46556    customer 
   TABLE DATA           o   COPY public.customer (id, first_name, second_name, login, password, email, type, subscription_end) FROM stdin;
    public          postgres    false    212   j4       D          0    46620    file 
   TABLE DATA           e   COPY public.file (id, project_id, filename, extension, size, file_path, update_datetime) FROM stdin;
    public          postgres    false    216   �5       ?          0    16386    goose_db_version 
   TABLE DATA           N   COPY public.goose_db_version (id, version_id, is_applied, tstamp) FROM stdin;
    public          postgres    false    211   7       E          0    46637    note 
   TABLE DATA           b   COPY public.note (id, project_id, title, content, update_datetime, deadline, overdue) FROM stdin;
    public          postgres    false    217   �7       B          0    46585    project 
   TABLE DATA           X   COPY public.project (id, owner_id, capacity, name, creation_date, admin_id) FROM stdin;
    public          postgres    false    214   8       C          0    46603    project_access 
   TABLE DATA           V   COPY public.project_access (id, project_id, customer_id, customer_access) FROM stdin;
    public          postgres    false    215   9       O           0    0    goose_db_version_id_seq    SEQUENCE SET     G   SELECT pg_catalog.setval('public.goose_db_version_id_seq', 523, true);
          public          postgres    false    210            �           2606    46584     administrator administrator_pkey 
   CONSTRAINT     ^   ALTER TABLE ONLY public.administrator
    ADD CONSTRAINT administrator_pkey PRIMARY KEY (id);
 J   ALTER TABLE ONLY public.administrator DROP CONSTRAINT administrator_pkey;
       public            postgres    false    213            �           2606    46573    customer customer_email_key 
   CONSTRAINT     W   ALTER TABLE ONLY public.customer
    ADD CONSTRAINT customer_email_key UNIQUE (email);
 E   ALTER TABLE ONLY public.customer DROP CONSTRAINT customer_email_key;
       public            postgres    false    212            �           2606    46571    customer customer_login_key 
   CONSTRAINT     W   ALTER TABLE ONLY public.customer
    ADD CONSTRAINT customer_login_key UNIQUE (login);
 E   ALTER TABLE ONLY public.customer DROP CONSTRAINT customer_login_key;
       public            postgres    false    212            �           2606    46569    customer customer_pkey 
   CONSTRAINT     T   ALTER TABLE ONLY public.customer
    ADD CONSTRAINT customer_pkey PRIMARY KEY (id);
 @   ALTER TABLE ONLY public.customer DROP CONSTRAINT customer_pkey;
       public            postgres    false    212            �           2606    46631    file file_pkey 
   CONSTRAINT     L   ALTER TABLE ONLY public.file
    ADD CONSTRAINT file_pkey PRIMARY KEY (id);
 8   ALTER TABLE ONLY public.file DROP CONSTRAINT file_pkey;
       public            postgres    false    216            �           2606    16392 &   goose_db_version goose_db_version_pkey 
   CONSTRAINT     d   ALTER TABLE ONLY public.goose_db_version
    ADD CONSTRAINT goose_db_version_pkey PRIMARY KEY (id);
 P   ALTER TABLE ONLY public.goose_db_version DROP CONSTRAINT goose_db_version_pkey;
       public            postgres    false    211            �           2606    46646    note note_pkey 
   CONSTRAINT     L   ALTER TABLE ONLY public.note
    ADD CONSTRAINT note_pkey PRIMARY KEY (id);
 8   ALTER TABLE ONLY public.note DROP CONSTRAINT note_pkey;
       public            postgres    false    217            �           2606    46609 "   project_access project_access_pkey 
   CONSTRAINT     `   ALTER TABLE ONLY public.project_access
    ADD CONSTRAINT project_access_pkey PRIMARY KEY (id);
 L   ALTER TABLE ONLY public.project_access DROP CONSTRAINT project_access_pkey;
       public            postgres    false    215            �           2606    46592    project project_pkey 
   CONSTRAINT     R   ALTER TABLE ONLY public.project
    ADD CONSTRAINT project_pkey PRIMARY KEY (id);
 >   ALTER TABLE ONLY public.project DROP CONSTRAINT project_pkey;
       public            postgres    false    214            �           2606    46632    file file_project_id_fkey    FK CONSTRAINT     �   ALTER TABLE ONLY public.file
    ADD CONSTRAINT file_project_id_fkey FOREIGN KEY (project_id) REFERENCES public.project(id) ON DELETE CASCADE;
 C   ALTER TABLE ONLY public.file DROP CONSTRAINT file_project_id_fkey;
       public          postgres    false    216    3238    214            �           2606    46647    note note_project_id_fkey    FK CONSTRAINT     �   ALTER TABLE ONLY public.note
    ADD CONSTRAINT note_project_id_fkey FOREIGN KEY (project_id) REFERENCES public.project(id) ON DELETE CASCADE;
 C   ALTER TABLE ONLY public.note DROP CONSTRAINT note_project_id_fkey;
       public          postgres    false    3238    217    214            �           2606    46615 .   project_access project_access_customer_id_fkey    FK CONSTRAINT     �   ALTER TABLE ONLY public.project_access
    ADD CONSTRAINT project_access_customer_id_fkey FOREIGN KEY (customer_id) REFERENCES public.customer(id) ON DELETE CASCADE;
 X   ALTER TABLE ONLY public.project_access DROP CONSTRAINT project_access_customer_id_fkey;
       public          postgres    false    215    212    3234            �           2606    46610 -   project_access project_access_project_id_fkey    FK CONSTRAINT     �   ALTER TABLE ONLY public.project_access
    ADD CONSTRAINT project_access_project_id_fkey FOREIGN KEY (project_id) REFERENCES public.project(id) ON DELETE CASCADE;
 W   ALTER TABLE ONLY public.project_access DROP CONSTRAINT project_access_project_id_fkey;
       public          postgres    false    214    3238    215            �           2606    46598    project project_admin_id_fkey    FK CONSTRAINT     �   ALTER TABLE ONLY public.project
    ADD CONSTRAINT project_admin_id_fkey FOREIGN KEY (admin_id) REFERENCES public.administrator(id) ON DELETE CASCADE;
 G   ALTER TABLE ONLY public.project DROP CONSTRAINT project_admin_id_fkey;
       public          postgres    false    214    3236    213            �           2606    46593    project project_owner_id_fkey    FK CONSTRAINT     �   ALTER TABLE ONLY public.project
    ADD CONSTRAINT project_owner_id_fkey FOREIGN KEY (owner_id) REFERENCES public.customer(id) ON DELETE CASCADE;
 G   ALTER TABLE ONLY public.project DROP CONSTRAINT project_owner_id_fkey;
       public          postgres    false    214    3234    212            A      x������ � �      @     x�]��N�@ ���+\����۹3-+RF��Ɛ���3P�C�T��u��|���Pr(P �5Z��jICQXo�"p\��b�<=ڞ�q�2[��iZ������|Y���~��,$�~�*�n!w9��ì�b�'����j���j4R�D�A��4�\Tj6�C'�R�Q2��tIjb���j��Cx<��sV����
z��y�t�x�{\�̥9�e��&�?O���;���1�ڢ���|5���A��(�~��[Y      D   �  x���Kj\1E�z�����_�Z�>��㦟i��x������@�Q&k��Q��$`L�A�D�{JD�D����X� �� ��t�Z��B&%Q��ՁI.A�%BA�&y��f�W��B�|���k���=�������}�5�-��o�j&�:�s����MZ�|�|?�E�2���!Hs �h����lÄ����:.˖�J�-z�Z���C�����^������o���4o�~%�5Q���ه耬�`jː3��lͮ栃�Mtu�����Ւ�&��6!���S���ǘ�G��04�xp$!�M:k�J�ێ%��^��,�����tS^��	>N5���tz�雹�i}D�����b|����5}+�	��N�Qûq�MM<      ?   x   x�u��!�7�b �&a�\׿.��jy���6��o�Ҫ�J\������q�#��v*���^��t,z3#���MO~�!N�=O>�
���'N�������y�O�������[�C>MD~D0=N      E   u   x�Eɱ� @����	Y��Y�H�����S$w�W��]&QH�R�/��3���j :�5�0� [_ު��vwo;,E����+��@�J�<k��'��|(��<r�o�z)1      B   �   x��лU�@��خ���kf�*�D�%	����!�������=GA*��B DHp��>eٳ� m�a�!�QEXZ'^���@\��?�?��BH�@�BtS���X�`^�+�Q3��~������,Q�`���6tv�'�s/��� �\���lhejļ�_�q�&��7�f�cv^m�����k��|�!�|��u�įwJ      C   _   x�̱�  �Zw!������GH���6(`��D!�k���FrYP#g�{�>?ru�����vCg��.��@�= ��O7�֖�8m:5���{� i��     