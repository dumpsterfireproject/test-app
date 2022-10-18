import React, { useState, useEffect } from 'react';
import LocalPharmacyIcon from '@mui/icons-material/LocalPharmacy';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import {Item} from './Item';

type Props = {
  items: Item[];
  selected: Item|undefined;
  setSelected: (val: Item) => void;
};

const Items: React.FC<Props> = ({
  items,
  selected,
  setSelected
}) => {

  const handleListItemClick = (
    event: React.MouseEvent<HTMLDivElement, MouseEvent>,
    item: Item,
  ) => {
    setSelected(item);
  };

  function isSelected(selected: Item|undefined, item: Item): boolean {
    if (selected === undefined) {
      return false;
    }
    return selected.id === item.id;
  }

  return (
    <React.Fragment>
      {items.map((item) => {
          return (
          <ListItemButton
            key = {item["id"]}
            selected={isSelected(selected, item)}
            onClick={(event) => handleListItemClick(event, item)}
          >
            <ListItemIcon>
              <LocalPharmacyIcon />
            </ListItemIcon>
            <ListItemText primary={item["description"]} />
          </ListItemButton>
          );
      })}
    </React.Fragment>
  );
};

export default Items;
