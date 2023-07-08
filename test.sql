CREATE TABLE `explorer_collection` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `name` varchar(1024) COLLATE utf8mb4_general_ci NOT NULL COMMENT '测试集名字',
  `description` text COLLATE utf8mb4_general_ci COMMENT '测试集记录描述',
  `source` tinyint NOT NULL DEFAULT '0' COMMENT '数据源，1：cloud-dev, 2：bam',
  `psm` varchar(128) COLLATE utf8mb4_general_ci NOT NULL COMMENT 'psm',
  `feature_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT 'bam feature id',
  `scope` int NOT NULL DEFAULT '0' COMMENT '0:private, 1:public',
  `authorized` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否授权给别人',
  `creator` varchar(128) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '创建者',
  `operator` varchar(128) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '更新者',
  `deleted` tinyint(1) NOT NULL COMMENT '软删除标志，1代表已删除',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_creator` (`creator`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='请求集合'