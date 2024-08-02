-- 给cs_task添加字段flag 用来标记是否是通过区域来创建的
ALTER TABLE cs_task
ADD COLUMN flag text NOT NULL DEFAULT '0';

COMMENT ON COLUMN cs_task.flag IS '是否同区域关联目标 1:是 0:否';