import * as React from 'react';
import Typography from '@mui/material/Typography';
import Title from './Title';
import {Item} from './Item';

type Props = {
  current: Item|undefined;
  items: Item[];
};

function currentSku(current: Item|undefined, items: Item[]): string {
  if (current !== undefined) {
    return current.sku;
  }
  if (items.length > 0) {
    return items[0].sku;
  }
  return "";
}

function currentDescription(current: Item|undefined, items: Item[]): string {
  if (current !== undefined) {
    return current.description;
  }
  if (items.length > 0) {
    return items[0].description;
  }
  return "";
}

export default function ItemDetails(props: Props) {
  return (
    <React.Fragment>
      <Title>Item Details</Title>
      <Typography component="p" variant="h4">
        {currentSku(props.current, props.items)}
      </Typography>
      <Typography color="text.secondary" sx={{ flex: 1 }}>
        {currentDescription(props.current, props.items)}
      </Typography>
    </React.Fragment>
  );
}