import sqlite3
import os

def transfer_matching_data(source_db_path, destination_db_path):
    try:
        # 确保源文件和目标文件存在
        if not os.path.exists(source_db_path):
            raise FileNotFoundError(f"源数据库文件 {source_db_path} 不存在！")
        if not os.path.exists(destination_db_path):
            raise FileNotFoundError(f"目标数据库文件 {destination_db_path} 不存在！")

        # 连接到源数据库
        source_conn = sqlite3.connect(source_db_path)
        source_cursor = source_conn.cursor()

        # 获取源数据库中的所有表名
        source_cursor.execute("SELECT name FROM sqlite_master WHERE type='table';")
        source_tables = [table[0] for table in source_cursor.fetchall()]

        if not source_tables:
            print("源数据库中没有表！")
            return

        # 连接到目标数据库
        destination_conn = sqlite3.connect(destination_db_path)
        destination_cursor = destination_conn.cursor()

        # 获取目标数据库中的所有表名
        destination_cursor.execute("SELECT name FROM sqlite_master WHERE type='table';")
        target_tables = [table[0] for table in destination_cursor.fetchall()]

        if not target_tables:
            print("目标数据库中没有表！")
            return

        # 定义需要传输的数据类型
        allowed_types = ["TEXT", "TIMESTAMP", "BLOB"]

        # 遍历源数据库中的每个表
        for source_table in source_tables:
            print(f"处理源表 {source_table}...")

            # 获取源表的字段及其类型
            source_cursor.execute(f"PRAGMA table_info({source_table});")
            source_fields = {col[1]: col[2] for col in source_cursor.fetchall() if col[2] in allowed_types}

            if not source_fields:
                print(f"源表 {source_table} 中没有匹配的字段类型（TEXT/TIMESTAMP/BLOB），跳过该表。")
                continue

            # 遍历目标数据库中的每个表，尝试匹配字段类型
            for target_table in target_tables:
                print(f"检查目标表 {target_table} 是否匹配...")

                # 获取目标表的字段及其类型
                destination_cursor.execute(f"PRAGMA table_info({target_table});")
                target_fields = {col[1]: col[2] for col in destination_cursor.fetchall() if col[2] in allowed_types}

                if not target_fields:
                    print(f"目标表 {target_table} 中没有匹配的字段类型，跳过该表。")
                    continue

                # 查找目标表中与源表字段类型相同的字段
                matching_fields = {}
                for source_field, source_type in source_fields.items():
                    for target_field, target_type in target_fields.items():
                        if source_type == target_type:  # 类型相同
                            matching_fields[source_field] = target_field
                            break

                if not matching_fields:
                    print(f"目标表 {target_table} 中没有与源表 {source_table} 匹配的字段类型。")
                    continue

                # 打印匹配的字段映射
                print(f"找到以下匹配的字段映射：{matching_fields}")

                # 从源表中获取数据
                select_fields = ", ".join(matching_fields.keys())
                source_cursor.execute(f"SELECT {select_fields} FROM {source_table};")
                rows = source_cursor.fetchall()

                if not rows:
                    print(f"源表 {source_table} 中没有数据，无需迁移。")
                    continue

                # 构建插入字段列表
                insert_fields = ", ".join(matching_fields.values())
                placeholders = ", ".join(["?"] * len(matching_fields))

                # 转换数据并插入目标表
                for row in rows:
                    # 使用匹配的字段映射重新组织数据
                    mapped_row = [row[list(matching_fields.keys()).index(source_field)] for source_field in matching_fields.keys()]
                    try:
                        destination_cursor.execute(f"INSERT INTO {target_table} ({insert_fields}) VALUES ({placeholders})", mapped_row)
                    except sqlite3.IntegrityError as e:
                        print(f"插入失败：{e}，跳过该记录。")
                        continue

                print(f"源表 {source_table} 的数据已成功迁移到目标表 {target_table}！共迁移 {len(rows)} 条记录。")

        # 提交更改并关闭连接
        destination_conn.commit()
        source_conn.close()
        destination_conn.close()

        print("数据迁移完成！")

    except sqlite3.Error as e:
        print(f"数据库错误: {e}")
    except Exception as e:
        print(f"发生了错误: {e}")

if __name__ == "__main__":
    # 设置源数据库和目标数据库的路径
    source_db_path = "memos_prod.db"
    destination_db_path = "noise.db"

    # 调用函数进行数据迁移
    transfer_matching_data(source_db_path, destination_db_path)
