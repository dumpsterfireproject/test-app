import React, {useState,useEffect} from 'react';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Title from './Title';
import {Item} from './Item';
import { Inventory } from './Inventory';

type Props = {
  current: Item|undefined;
};

export default function InventoryDetails(props: Props) {
  const [currentInventory, setCurrentInventory] = useState<Inventory[]>([]);
  useEffect(() => {
    if (props.current !== undefined) {
      fetch(
        `http://localhost:8080/api/inventory?itemID=${props.current.id}`,
        {
          method: 'GET', 
          headers: {  
            'Authorization': 'Bearer 1234'  
          }
        },
      )
      .then((response) => response.json())
      .then((data) => {
          console.log(data);
          setCurrentInventory(data);
      })
      .catch((err) => {
          console.log(err.message);
      });
    }
 }, [props]);

  return (
    <React.Fragment>
      <Title>Inventory</Title>
      <Table size="small">
        <TableHead>
          <TableRow>
            <TableCell>Location</TableCell>
            <TableCell>Status</TableCell>
            <TableCell align="right">Quantity</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {currentInventory.map((inventory) => (
            <TableRow key={inventory.id}>
              <TableCell>{inventory.location}</TableCell>
              <TableCell>{inventory.status}</TableCell>
              <TableCell align="right">{inventory.quantity}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </React.Fragment>
  );
}