create table if not exists sys_db_changelog
(
    script_name varchar(300) not null
        primary key,
    create_ts   timestamp default CURRENT_TIMESTAMP,
    is_init     integer   default 0
);

alter table sys_db_changelog
    owner to postgres;

create table if not exists sys_server
(
    id         uuid not null
        primary key,
    create_ts  timestamp,
    created_by varchar(50),
    update_ts  timestamp,
    updated_by varchar(50),
    name       varchar(255),
    is_running boolean,
    data       text
);

alter table sys_server
    owner to postgres;

create unique index if not exists idx_sys_server_uniq_name
    on sys_server (name);

create table if not exists sys_config
(
    id         uuid         not null
        primary key,
    create_ts  timestamp,
    created_by varchar(50),
    version    integer,
    update_ts  timestamp,
    updated_by varchar(50),
    name       varchar(255) not null,
    value_     text         not null
);

alter table sys_config
    owner to postgres;

create unique index if not exists idx_sys_config_uniq_name
    on sys_config (name);

create table if not exists sys_file
(
    id          uuid         not null
        primary key,
    create_ts   timestamp,
    created_by  varchar(50),
    version     integer,
    update_ts   timestamp,
    updated_by  varchar(50),
    delete_ts   timestamp,
    deleted_by  varchar(50),
    name        varchar(500) not null,
    ext         varchar(20),
    file_size   bigint,
    create_date timestamp
);

alter table sys_file
    owner to postgres;

create table if not exists sys_lock_config
(
    id          uuid not null
        primary key,
    create_ts   timestamp,
    created_by  varchar(50),
    name        varchar(100),
    timeout_sec integer
);

alter table sys_lock_config
    owner to postgres;

create table if not exists sys_entity_statistics
(
    id                        uuid not null
        primary key,
    create_ts                 timestamp,
    created_by                varchar(50),
    update_ts                 timestamp,
    updated_by                varchar(50),
    name                      varchar(50),
    instance_count            bigint,
    fetch_ui                  integer,
    max_fetch_ui              integer,
    lazy_collection_threshold integer,
    lookup_screen_threshold   integer
);

alter table sys_entity_statistics
    owner to postgres;

create unique index if not exists idx_sys_entity_statistics_uniq_name
    on sys_entity_statistics (name);

create table if not exists sys_scheduled_task
(
    id                uuid not null
        primary key,
    create_ts         timestamp,
    created_by        varchar(50),
    update_ts         timestamp,
    updated_by        varchar(50),
    delete_ts         timestamp,
    deleted_by        varchar(50),
    defined_by        varchar(1) default 'B'::character varying,
    class_name        varchar(500),
    script_name       varchar(500),
    bean_name         varchar(50),
    method_name       varchar(50),
    method_params     varchar(1000),
    user_name         varchar(50),
    is_singleton      boolean,
    is_active         boolean,
    period_           integer,
    timeout           integer,
    start_date        timestamp,
    time_frame        integer,
    start_delay       integer,
    permitted_servers varchar(4096),
    log_start         boolean,
    log_finish        boolean,
    last_start_time   timestamp with time zone,
    last_start_server varchar(512),
    description       varchar(1000),
    cron              varchar(100),
    scheduling_type   varchar(1) default 'P'::character varying
);

alter table sys_scheduled_task
    owner to postgres;

create table if not exists sys_scheduled_execution
(
    id          uuid not null
        primary key,
    create_ts   timestamp,
    created_by  varchar(50),
    task_id     uuid
        constraint sys_scheduled_execution_task
            references sys_scheduled_task,
    server      varchar(512),
    start_time  timestamp with time zone,
    finish_time timestamp with time zone,
    result      text
);

alter table sys_scheduled_execution
    owner to postgres;

create index if not exists idx_sys_scheduled_execution_task_start_time
    on sys_scheduled_execution (task_id, start_time);

create index if not exists idx_sys_scheduled_execution_task_finish_time
    on sys_scheduled_execution (task_id, finish_time);

create table if not exists sec_role
(
    id              uuid         not null
        primary key,
    create_ts       timestamp,
    created_by      varchar(50),
    version         integer,
    update_ts       timestamp,
    updated_by      varchar(50),
    delete_ts       timestamp,
    deleted_by      varchar(50),
    name            varchar(255) not null,
    loc_name        varchar(255),
    description     varchar(1000),
    is_default_role boolean,
    role_type       integer
);

alter table sec_role
    owner to postgres;

create unique index if not exists idx_sec_role_uniq_name
    on sec_role (name)
    where (delete_ts IS NULL);

create table if not exists sec_group
(
    id         uuid         not null
        primary key,
    create_ts  timestamp,
    created_by varchar(50),
    version    integer,
    update_ts  timestamp,
    updated_by varchar(50),
    delete_ts  timestamp,
    deleted_by varchar(50),
    name       varchar(255) not null,
    parent_id  uuid
        constraint sec_group_parent
            references sec_group
);

alter table sec_group
    owner to postgres;

create unique index if not exists idx_sec_group_uniq_name
    on sec_group (name)
    where (delete_ts IS NULL);

create table if not exists sec_group_hierarchy
(
    id              uuid not null
        primary key,
    create_ts       timestamp,
    created_by      varchar(50),
    group_id        uuid
        constraint sec_group_hierarchy_group
            references sec_group,
    parent_id       uuid
        constraint sec_group_hierarchy_parent
            references sec_group,
    hierarchy_level integer
);

alter table sec_group_hierarchy
    owner to postgres;

create table if not exists sec_user
(
    id                       uuid        not null
        primary key,
    create_ts                timestamp,
    created_by               text,
    version                  bigint,
    update_ts                timestamp,
    updated_by               text,
    delete_ts                timestamp,
    deleted_by               text,
    login                    text,
    login_lc                 varchar(50) not null,
    password                 varchar(255),
    password_encryption      varchar(50),
    name                     text,
    first_name               text,
    last_name                text,
    middle_name              text,
    position_                varchar(255),
    email                    text,
    language_                varchar(20),
    time_zone                varchar(50),
    time_zone_auto           boolean,
    active                   boolean,
    group_id                 uuid        not null
        constraint sec_user_group
            references sec_group,
    ip_mask                  varchar(200),
    change_password_at_logon boolean
);

alter table sec_user
    owner to postgres;

create unique index if not exists idx_sec_user_uniq_login
    on sec_user (login_lc)
    where (delete_ts IS NULL);

create table if not exists sec_user_role
(
    id         uuid not null
        primary key,
    create_ts  timestamp,
    created_by varchar(50),
    version    integer,
    update_ts  timestamp,
    updated_by varchar(50),
    delete_ts  timestamp,
    deleted_by varchar(50),
    user_id    uuid
        constraint sec_user_role_profile
            references sec_user,
    role_id    uuid
        constraint sec_user_role_role
            references sec_role
);

alter table sec_user_role
    owner to postgres;

create unique index if not exists idx_sec_user_role_uniq_role
    on sec_user_role (user_id, role_id)
    where (delete_ts IS NULL);

create table if not exists sec_permission
(
    id              uuid not null
        primary key,
    create_ts       timestamp,
    created_by      varchar(50),
    version         integer,
    update_ts       timestamp,
    updated_by      varchar(50),
    delete_ts       timestamp,
    deleted_by      varchar(50),
    permission_type integer,
    target          varchar(100),
    value_          integer,
    role_id         uuid
        constraint sec_permission_role
            references sec_role
);

alter table sec_permission
    owner to postgres;

create unique index if not exists idx_sec_permission_unique
    on sec_permission (role_id, permission_type, target)
    where (delete_ts IS NULL);

create table if not exists sec_constraint
(
    id             uuid         not null
        primary key,
    create_ts      timestamp,
    created_by     varchar(50),
    version        integer,
    update_ts      timestamp,
    updated_by     varchar(50),
    delete_ts      timestamp,
    deleted_by     varchar(50),
    code           varchar(255),
    check_type     varchar(50) default 'db'::character varying,
    operation_type varchar(50) default 'read'::character varying,
    entity_name    varchar(255) not null,
    join_clause    varchar(500),
    where_clause   varchar(1000),
    groovy_script  text,
    filter_xml     text,
    is_active      boolean     default true,
    group_id       uuid
        constraint sec_constraint_group
            references sec_group
);

alter table sec_constraint
    owner to postgres;

create index if not exists idx_sec_constraint_group
    on sec_constraint (group_id);

create table if not exists sec_localized_constraint_msg
(
    id             uuid         not null
        primary key,
    create_ts      timestamp,
    created_by     varchar(50),
    version        integer,
    update_ts      timestamp,
    updated_by     varchar(50),
    delete_ts      timestamp,
    deleted_by     varchar(50),
    entity_name    varchar(255) not null,
    operation_type varchar(50)  not null,
    values_        text
);

alter table sec_localized_constraint_msg
    owner to postgres;

create unique index if not exists idx_sec_loc_cnstrnt_msg_unique
    on sec_localized_constraint_msg (entity_name, operation_type)
    where (delete_ts IS NULL);

create table if not exists sec_session_attr
(
    id         uuid not null
        primary key,
    create_ts  timestamp,
    created_by varchar(50),
    version    integer,
    update_ts  timestamp,
    updated_by varchar(50),
    delete_ts  timestamp,
    deleted_by varchar(50),
    name       varchar(50),
    str_value  varchar(1000),
    datatype   varchar(20),
    group_id   uuid
        constraint sec_session_attr_group
            references sec_group
);

alter table sec_session_attr
    owner to postgres;

create index if not exists idx_sec_session_attr_group
    on sec_session_attr (group_id);

create table if not exists sec_user_setting
(
    id          uuid not null
        primary key,
    create_ts   timestamp,
    created_by  varchar(50),
    user_id     uuid
        constraint sec_user_setting_user
            references sec_user,
    client_type char,
    name        varchar(255),
    value_      text,
    constraint sec_user_setting_uniq
        unique (user_id, name, client_type)
);

alter table sec_user_setting
    owner to postgres;

create table if not exists sec_user_substitution
(
    id                  uuid not null
        primary key,
    create_ts           timestamp,
    created_by          varchar(50),
    version             integer,
    update_ts           timestamp,
    updated_by          varchar(50),
    delete_ts           timestamp,
    deleted_by          varchar(50),
    user_id             uuid not null
        constraint fk_sec_user_substitution_user
            references sec_user,
    substituted_user_id uuid not null
        constraint fk_sec_user_substitution_substituted_user
            references sec_user,
    start_date          timestamp,
    end_date            timestamp
);

alter table sec_user_substitution
    owner to postgres;

create index if not exists idx_sec_user_substitution_user
    on sec_user_substitution (user_id);

create table if not exists sec_logged_entity
(
    id         uuid not null
        primary key,
    create_ts  timestamp,
    created_by varchar(50),
    name       varchar(100)
        constraint sec_logged_entity_uniq_name
            unique,
    auto       boolean,
    manual     boolean
);

alter table sec_logged_entity
    owner to postgres;

create table if not exists sec_logged_attr
(
    id         uuid not null
        primary key,
    create_ts  timestamp,
    created_by varchar(50),
    entity_id  uuid
        constraint fk_sec_logged_attr_entity
            references sec_logged_entity,
    name       varchar(50),
    constraint sec_logged_attr_uniq_name
        unique (entity_id, name)
);

alter table sec_logged_attr
    owner to postgres;

create index if not exists idx_sec_logged_attr_entity
    on sec_logged_attr (entity_id);

create table if not exists sec_entity_log
(
    id                   uuid not null
        primary key,
    create_ts            timestamp,
    created_by           varchar(50),
    event_ts             timestamp,
    user_id              uuid
        constraint fk_sec_entity_log_user
            references sec_user,
    change_type          char,
    entity               varchar(100),
    entity_instance_name varchar(1000),
    entity_id            uuid,
    string_entity_id     varchar(255),
    int_entity_id        integer,
    long_entity_id       bigint,
    changes              text
);

alter table sec_entity_log
    owner to postgres;

create index if not exists idx_sec_entity_log_entity_id
    on sec_entity_log (entity_id);

create index if not exists idx_sec_entity_log_sentity_id
    on sec_entity_log (string_entity_id);

create index if not exists idx_sec_entity_log_ientity_id
    on sec_entity_log (int_entity_id);

create index if not exists idx_sec_entity_log_lentity_id
    on sec_entity_log (long_entity_id);

create table if not exists sec_filter
(
    id             uuid not null
        primary key,
    create_ts      timestamp,
    created_by     varchar(50),
    version        integer,
    update_ts      timestamp,
    updated_by     varchar(50),
    delete_ts      timestamp,
    deleted_by     varchar(50),
    component      varchar(200),
    name           varchar(255),
    code           varchar(200),
    xml            text,
    user_id        uuid
        constraint fk_sec_filter_user
            references sec_user,
    global_default boolean
);

alter table sec_filter
    owner to postgres;

create index if not exists idx_sec_filter_component_user
    on sec_filter (component, user_id);

create table if not exists sys_folder
(
    id          uuid not null
        primary key,
    create_ts   timestamp,
    created_by  varchar(50),
    version     integer,
    update_ts   timestamp,
    updated_by  varchar(50),
    delete_ts   timestamp,
    deleted_by  varchar(50),
    folder_type char,
    parent_id   uuid
        constraint fk_sys_folder_parent
            references sys_folder,
    name        varchar(100),
    tab_name    varchar(100),
    sort_order  integer
);

alter table sys_folder
    owner to postgres;

create table if not exists sys_app_folder
(
    folder_id         uuid not null
        primary key
        constraint fk_sys_app_folder_folder
            references sys_folder,
    filter_component  varchar(200),
    filter_xml        varchar(7000),
    visibility_script text,
    quantity_script   text,
    apply_default     boolean
);

alter table sys_app_folder
    owner to postgres;

create table if not exists sec_presentation
(
    id           uuid not null
        primary key,
    create_ts    timestamp,
    created_by   varchar(50),
    update_ts    timestamp,
    updated_by   varchar(50),
    component    varchar(200),
    name         varchar(255),
    xml          varchar(7000),
    user_id      uuid
        constraint sec_presentation_user
            references sec_user,
    is_auto_save boolean
);

alter table sec_presentation
    owner to postgres;

create index if not exists idx_sec_presentation_component_user
    on sec_presentation (component, user_id);

create table if not exists sec_search_folder
(
    folder_id        uuid not null
        primary key
        constraint fk_sec_search_folder_folder
            references sys_folder,
    filter_component varchar(200),
    filter_xml       varchar(7000),
    user_id          uuid
        constraint fk_sec_search_folder_user
            references sec_user,
    presentation_id  uuid
        constraint fk_sec_search_folder_presentation
            references sec_presentation
            on delete set null,
    apply_default    boolean,
    is_set           boolean,
    entity_type      varchar(50)
);

alter table sec_search_folder
    owner to postgres;

create index if not exists idx_sec_search_folder_user
    on sec_search_folder (user_id);

create table if not exists sys_fts_queue
(
    id               uuid not null
        primary key,
    create_ts        timestamp,
    created_by       varchar(50),
    entity_id        uuid,
    string_entity_id varchar(255),
    int_entity_id    integer,
    long_entity_id   bigint,
    entity_name      varchar(200),
    change_type      char,
    source_host      varchar(255),
    indexing_host    varchar(255),
    fake             boolean
);

alter table sys_fts_queue
    owner to postgres;

create index if not exists idx_sys_fts_queue_idxhost_crts
    on sys_fts_queue (indexing_host, create_ts);

create table if not exists sec_screen_history
(
    id                  uuid not null
        primary key,
    create_ts           timestamp,
    created_by          varchar(50),
    user_id             uuid
        constraint fk_sec_history_user
            references sec_user,
    caption             varchar(255),
    url                 text,
    entity_id           uuid,
    string_entity_id    varchar(255),
    int_entity_id       integer,
    long_entity_id      bigint,
    substituted_user_id uuid
        constraint fk_sec_history_substituted_user
            references sec_user
);

alter table sec_screen_history
    owner to postgres;

create index if not exists idx_sec_screen_history_user
    on sec_screen_history (user_id);

create index if not exists idx_sec_screen_hist_sub_user
    on sec_screen_history (substituted_user_id);

create index if not exists idx_sec_screen_history_entity_id
    on sec_screen_history (entity_id);

create index if not exists idx_sec_screen_history_sentity_id
    on sec_screen_history (string_entity_id);

create index if not exists idx_sec_screen_history_ientity_id
    on sec_screen_history (int_entity_id);

create index if not exists idx_sec_screen_history_lentity_id
    on sec_screen_history (long_entity_id);

create table if not exists sys_sending_message
(
    id                   uuid not null
        primary key,
    create_ts            timestamp,
    created_by           varchar(50),
    version              integer,
    update_ts            timestamp with time zone,
    updated_by           varchar(50),
    delete_ts            timestamp,
    deleted_by           varchar(50),
    address_to           text,
    address_cc           text,
    address_bcc          text,
    address_from         varchar(100),
    caption              varchar(500),
    email_headers        varchar(500),
    content_text         text,
    content_text_file_id uuid
        constraint fk_sys_sending_message_content_file
            references sys_file,
    deadline             timestamp with time zone,
    status               integer,
    date_sent            timestamp,
    attempts_count       integer,
    attempts_made        integer,
    attachments_name     text,
    body_content_type    varchar(50)
);

alter table sys_sending_message
    owner to postgres;

create index if not exists idx_sys_sending_message_status
    on sys_sending_message (status);

create index if not exists idx_sys_sending_message_date_sent
    on sys_sending_message (date_sent);

create index if not exists idx_sys_sending_message_update_ts
    on sys_sending_message (update_ts);

create table if not exists sys_sending_attachment
(
    id              uuid not null
        primary key,
    create_ts       timestamp,
    created_by      varchar(50),
    version         integer,
    update_ts       timestamp,
    updated_by      varchar(50),
    delete_ts       timestamp,
    deleted_by      varchar(50),
    message_id      uuid
        constraint fk_sys_sending_attachment_sending_message
            references sys_sending_message,
    content         bytea,
    content_file_id uuid
        constraint fk_sys_sending_attachment_content_file
            references sys_file,
    content_id      varchar(50),
    name            varchar(500),
    disposition     varchar(50),
    text_encoding   varchar(50)
);

alter table sys_sending_attachment
    owner to postgres;

create index if not exists sys_sending_attachment_message_idx
    on sys_sending_attachment (message_id);

create table if not exists sys_entity_snapshot
(
    id                uuid        not null
        primary key,
    create_ts         timestamp,
    created_by        varchar(50),
    entity_meta_class varchar(50) not null,
    entity_id         uuid,
    string_entity_id  varchar(255),
    int_entity_id     integer,
    long_entity_id    bigint,
    author_id         uuid        not null
        constraint fk_sys_entity_snapshot_author_id
            references sec_user,
    view_xml          text        not null,
    snapshot_xml      text        not null,
    snapshot_date     timestamp   not null
);

alter table sys_entity_snapshot
    owner to postgres;

create index if not exists idx_sys_entity_snapshot_entity_id
    on sys_entity_snapshot (entity_id);

create index if not exists idx_sys_entity_snapshot_sentity_id
    on sys_entity_snapshot (string_entity_id);

create index if not exists idx_sys_entity_snapshot_ientity_id
    on sys_entity_snapshot (int_entity_id);

create index if not exists idx_sys_entity_snapshot_lentity_id
    on sys_entity_snapshot (long_entity_id);

create table if not exists sys_category
(
    id            uuid         not null
        primary key,
    create_ts     timestamp,
    created_by    varchar(50),
    version       integer,
    update_ts     timestamp,
    updated_by    varchar(50),
    delete_ts     timestamp,
    deleted_by    varchar(50),
    name          varchar(255) not null,
    special       varchar(50),
    entity_type   varchar(100) not null,
    is_default    boolean,
    discriminator integer,
    locale_names  varchar(1000)
);

alter table sys_category
    owner to postgres;

create unique index if not exists idx_sys_category_uniq_name_entity_type
    on sys_category (name, entity_type)
    where (delete_ts IS NULL);

create table if not exists sys_category_attr
(
    id                           uuid         not null
        primary key,
    create_ts                    timestamp,
    created_by                   varchar(50),
    version                      integer,
    update_ts                    timestamp,
    updated_by                   varchar(50),
    delete_ts                    timestamp,
    deleted_by                   varchar(50),
    category_entity_type         varchar(4000),
    name                         varchar(255),
    code                         varchar(100) not null,
    description                  varchar(1000),
    category_id                  uuid         not null
        constraint sys_category_attr_category_id
            references sys_category,
    entity_class                 varchar(500),
    data_type                    varchar(200),
    default_string               varchar,
    default_int                  integer,
    default_double               numeric(36, 6),
    default_decimal              numeric(36, 10),
    default_date                 timestamp,
    default_date_wo_time         date,
    default_date_is_current      boolean,
    default_boolean              boolean,
    default_entity_value         uuid,
    default_str_entity_value     varchar(255),
    default_int_entity_value     integer,
    default_long_entity_value    bigint,
    enumeration                  varchar(500),
    order_no                     integer,
    screen                       varchar(255),
    required                     boolean,
    lookup                       boolean,
    target_screens               varchar(4000),
    width                        varchar(20),
    rows_count                   integer,
    is_collection                boolean,
    join_clause                  varchar(4000),
    where_clause                 varchar(4000),
    filter_xml                   text,
    locale_names                 varchar(1000),
    locale_descriptions          varchar(4000),
    enumeration_locales          varchar(5000),
    attribute_configuration_json text
);

alter table sys_category_attr
    owner to postgres;

create index if not exists idx_sys_category_attr_category
    on sys_category_attr (category_id);

create unique index if not exists idx_cat_attr_ent_type_and_code
    on sys_category_attr (category_entity_type, code)
    where (delete_ts IS NULL);

create table if not exists sys_attr_value
(
    id                  uuid         not null
        primary key,
    create_ts           timestamp,
    created_by          varchar(50),
    version             integer,
    update_ts           timestamp,
    updated_by          varchar(50),
    delete_ts           timestamp,
    deleted_by          varchar(50),
    category_attr_id    uuid         not null
        constraint sys_attr_value_category_attr_id
            references sys_category_attr,
    code                varchar(100) not null,
    entity_id           uuid,
    string_entity_id    varchar(255),
    int_entity_id       integer,
    long_entity_id      bigint,
    string_value        varchar,
    integer_value       integer,
    double_value        numeric(36, 6),
    decimal_value       numeric(36, 10),
    date_value          timestamp,
    date_wo_time_value  date,
    boolean_value       boolean,
    entity_value        uuid,
    string_entity_value varchar(255),
    int_entity_value    integer,
    long_entity_value   bigint,
    parent_id           uuid
        constraint sys_attr_value_attr_value_parent_id
            references sys_attr_value
);

alter table sys_attr_value
    owner to postgres;

create index if not exists idx_sys_attr_value_entity
    on sys_attr_value (entity_id);

create index if not exists idx_sys_attr_value_sentity
    on sys_attr_value (string_entity_id);

create index if not exists idx_sys_attr_value_ientity
    on sys_attr_value (int_entity_id);

create index if not exists idx_sys_attr_value_lentity
    on sys_attr_value (long_entity_id);

create table if not exists sys_jmx_instance
(
    id         uuid         not null
        primary key,
    create_ts  timestamp,
    created_by varchar(50),
    version    integer,
    update_ts  timestamp,
    updated_by varchar(50),
    delete_ts  timestamp,
    deleted_by varchar(50),
    node_name  varchar(255),
    address    varchar(500) not null,
    login      varchar(50)  not null,
    password   varchar(255) not null
);

alter table sys_jmx_instance
    owner to postgres;

create sequence if not exists sys_query_result_seq;

alter sequence sys_query_result_seq owner to postgres;

create table if not exists sys_query_result
(
    id               bigint default nextval('sys_query_result_seq'::regclass) not null
        primary key,
    session_id       uuid                                                     not null,
    query_key        integer                                                  not null,
    entity_id        uuid,
    string_entity_id varchar(255),
    int_entity_id    integer,
    long_entity_id   bigint
);

alter table sys_query_result
    owner to postgres;

create index if not exists idx_sys_query_result_entity_session_key
    on sys_query_result (entity_id, session_id, query_key);

create index if not exists idx_sys_query_result_sentity_session_key
    on sys_query_result (string_entity_id, session_id, query_key);

create index if not exists idx_sys_query_result_ientity_session_key
    on sys_query_result (int_entity_id, session_id, query_key);

create index if not exists idx_sys_query_result_lentity_session_key
    on sys_query_result (long_entity_id, session_id, query_key);

create index if not exists idx_sys_query_result_session_key
    on sys_query_result (session_id, query_key);

create table if not exists sec_remember_me
(
    id         uuid        not null
        primary key,
    create_ts  timestamp,
    created_by varchar(50),
    version    integer,
    user_id    uuid        not null
        constraint fk_sec_remember_me_user
            references sec_user,
    token      varchar(32) not null
);

alter table sec_remember_me
    owner to postgres;

create index if not exists idx_sec_remember_me_user
    on sec_remember_me (user_id);

create index if not exists idx_sec_remember_me_token
    on sec_remember_me (token);

create table if not exists sec_session_log
(
    id                  uuid    not null
        primary key,
    version             integer not null,
    create_ts           timestamp,
    created_by          varchar(50),
    update_ts           timestamp,
    updated_by          varchar(50),
    delete_ts           timestamp,
    deleted_by          varchar(50),
    session_id          uuid    not null,
    user_id             uuid    not null
        constraint fk_sec_session_log_user
            references sec_user,
    substituted_user_id uuid
        constraint fk_sec_session_log_subuser
            references sec_user,
    user_data           text,
    last_action         integer not null,
    client_info         varchar(512),
    client_type         varchar(10),
    address             varchar(255),
    started_ts          timestamp,
    finished_ts         timestamp,
    server_id           varchar(128)
);

alter table sec_session_log
    owner to postgres;

create index if not exists idx_sec_session_log_user
    on sec_session_log (user_id);

create index if not exists idx_sec_session_log_subuser
    on sec_session_log (substituted_user_id);

create index if not exists idx_sec_session_log_session
    on sec_session_log (session_id);

create index if not exists idx_session_log_started_ts
    on sec_session_log (started_ts desc);

create table if not exists sys_access_token
(
    id                   uuid not null
        primary key,
    create_ts            timestamp,
    token_value          varchar(255),
    token_bytes          bytea,
    authentication_key   varchar(255),
    authentication_bytes bytea,
    expiry               timestamp,
    user_login           varchar(50),
    locale               varchar(200),
    refresh_token_value  varchar(255)
);

alter table sys_access_token
    owner to postgres;

create table if not exists sys_refresh_token
(
    id                   uuid not null
        primary key,
    create_ts            timestamp,
    token_value          varchar(255),
    token_bytes          bytea,
    authentication_bytes bytea,
    expiry               timestamp,
    user_login           varchar(50)
);

alter table sys_refresh_token
    owner to postgres;

insert into SEC_GROUP (ID, CREATE_TS, VERSION, NAME, PARENT_ID)
values ('0fa2b1a5-1d68-4d69-9fbd-dff348347f93', now(), 0, 'Company', null) on conflict do nothing;

insert into SEC_USER (ID, CREATE_TS, VERSION, LOGIN, LOGIN_LC, PASSWORD, PASSWORD_ENCRYPTION, NAME, GROUP_ID, ACTIVE)
values ('60885987-1b61-4247-94c7-dff348347f93', now(), 0, 'admin', 'admin',
        '$2a$10$vQx8b8B7jzZ0rQmtuK4YDOKp7nkmUCFjPx6DMT.voPtetNHFOsaOu', 'bcrypt',
        'Administrator', '0fa2b1a5-1d68-4d69-9fbd-dff348347f93', true) on conflict do nothing;

insert into SEC_USER (ID, CREATE_TS, VERSION, LOGIN, LOGIN_LC, PASSWORD, NAME, GROUP_ID, ACTIVE)
values ('a405db59-e674-4f63-8afe-269dda788fe8', now(), 0, 'anonymous', 'anonymous', null,
        'Anonymous', '0fa2b1a5-1d68-4d69-9fbd-dff348347f93', true) on conflict do nothing;

insert into SEC_ROLE (ID, CREATE_TS, VERSION, NAME, ROLE_TYPE)
values ('0c018061-b26f-4de2-a5be-dff348347f93', now(), 0, 'Administrators', 10) on conflict do nothing;

insert into SEC_ROLE (ID, CREATE_TS, VERSION, NAME, ROLE_TYPE)
values ('cd541dd4-eeb7-cd5b-847e-d32236552fa9', now(), 0, 'Anonymous', 30) on conflict do nothing;

insert into SEC_USER_ROLE (ID, CREATE_TS, VERSION, USER_ID, ROLE_ID)
values ('c838be0a-96d0-4ef4-a7c0-dff348347f93', now(), 0, '60885987-1b61-4247-94c7-dff348347f93', '0c018061-b26f-4de2-a5be-dff348347f93') on conflict do nothing;

insert into SEC_USER_ROLE (ID, CREATE_TS, VERSION, USER_ID, ROLE_ID)
values ('f01fb532-c2f0-dc18-b86c-450cf8a8d8c5', now(), 0, 'a405db59-e674-4f63-8afe-269dda788fe8', 'cd541dd4-eeb7-cd5b-847e-d32236552fa9') on conflict do nothing;

insert into SEC_FILTER (ID, CREATE_TS, CREATED_BY, VERSION, COMPONENT, NAME, XML, USER_ID, GLOBAL_DEFAULT)
values ('b61d18cb-e79a-46f3-b16d-eaf4aebb10dd', now(), 'admin', 0, '[sec$User.browse].genericFilter', 'Search by role',
        '<?xml version="1.0" encoding="UTF-8"?><filter><and><c name="UrMxpkfMGn" class="com.haulmont.cuba.security.entity.Role" type="CUSTOM" locCaption="Role" entityAlias="u" join="join u.userRoles ur">ur.role.id = :component$genericFilter.UrMxpkfMGn32565<param name="component$genericFilter.UrMxpkfMGn32565">NULL</param></c></and></filter>',
        '60885987-1b61-4247-94c7-dff348347f93', false) on conflict do nothing;
