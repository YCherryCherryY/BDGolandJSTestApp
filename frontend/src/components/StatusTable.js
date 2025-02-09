import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { Table } from 'antd';

const StatusTable = () => {
  const [statuses, setStatuses] = useState([]);

  useEffect(() => {
    // Запрос данных с Backend API
    axios.get('http://localhost:8080/statuses')
      .then(response => {
        setStatuses(response.data);
      })
      .catch(error => {
        console.error('Error fetching statuses:', error);
      });
  }, []);

  // Определение колонок таблицы
  const columns = [
    {
      title: 'IP',
      dataIndex: 'ip',
      key: 'ip',
    },
    {
      title: 'Status',
      dataIndex: 'status',
      key: 'status',
    },
    {
      title: 'Check Time',
      dataIndex: 'time',
      key: 'time',
      render: (text) => new Date(text).toLocaleString(),
    },
    {
      title: 'Last Success Time',
      dataIndex: 'succsesTime',
      key: 'succsesTime',
      render: (text) => new Date(text).toLocaleString(),
    },
  ];

  return (
    <Table
      dataSource={statuses}
      columns={columns}
      rowKey="id"
      pagination={false}
    />
  );
};

export default StatusTable;