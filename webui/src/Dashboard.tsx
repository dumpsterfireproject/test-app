import React, {useState,useEffect} from 'react';
import { styled, createTheme, ThemeProvider } from '@mui/material/styles';
import Box from '@mui/material/Box';
import ChevronLeftIcon from '@mui/icons-material/ChevronLeft';
import Container from '@mui/material/Container';
import CssBaseline from '@mui/material/CssBaseline';
import Divider from '@mui/material/Divider';
import Grid from '@mui/material/Grid';
import IconButton from '@mui/material/IconButton';
import List from '@mui/material/List';
import MenuIcon from '@mui/icons-material/Menu';
import MuiAppBar, { AppBarProps as MuiAppBarProps } from '@mui/material/AppBar';
import MuiDrawer from '@mui/material/Drawer';
import Paper from '@mui/material/Paper';
import Toolbar from '@mui/material/Toolbar';
import Typography from '@mui/material/Typography';
import CreateIcon from '@mui/icons-material/Create';
import DeleteIcon from '@mui/icons-material/Delete';
import {Item} from './Item';
import Items from './Items';
import InventoryDetails from './InventoryDetails';
import ItemDetails from './ItemDetails';
import { stringify } from 'querystring';

import Button from '@mui/material/Button';
import TextField from '@mui/material/TextField';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogTitle from '@mui/material/DialogTitle';
import Select, { SelectChangeEvent } from '@mui/material/Select';
import InputLabel from '@mui/material/InputLabel';
import MenuItem from '@mui/material/MenuItem';

const drawerWidth: number = 240;

interface AppBarProps extends MuiAppBarProps {
  open?: boolean;
}

const AppBar = styled(MuiAppBar, {
  shouldForwardProp: (prop) => prop !== 'open',
})<AppBarProps>(({ theme, open }) => ({
  zIndex: theme.zIndex.drawer + 1,
  transition: theme.transitions.create(['width', 'margin'], {
    easing: theme.transitions.easing.sharp,
    duration: theme.transitions.duration.leavingScreen,
  }),
  ...(open && {
    marginLeft: drawerWidth,
    width: `calc(100% - ${drawerWidth}px)`,
    transition: theme.transitions.create(['width', 'margin'], {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.enteringScreen,
    }),
  }),
}));

const Drawer = styled(MuiDrawer, { shouldForwardProp: (prop) => prop !== 'open' })(
    ({ theme, open }) => ({
      '& .MuiDrawer-paper': {
        position: 'relative',
        whiteSpace: 'nowrap',
        width: drawerWidth,
        transition: theme.transitions.create('width', {
          easing: theme.transitions.easing.sharp,
          duration: theme.transitions.duration.enteringScreen,
        }),
        boxSizing: 'border-box',
        ...(!open && {
          overflowX: 'hidden',
          transition: theme.transitions.create('width', {
            easing: theme.transitions.easing.sharp,
            duration: theme.transitions.duration.leavingScreen,
          }),
          width: theme.spacing(7),
          [theme.breakpoints.up('sm')]: {
            width: theme.spacing(9),
          },
        }),
      },
    }),
  );
  
  const mdTheme = createTheme();
  
  function DashboardContent() {
    const [drawerOpen, setDrawerOpen] = React.useState(true);
    const toggleDrawer = () => {
      setDrawerOpen(!drawerOpen);
    };

    const [items, setItems] = useState<Item[]>([]);
    const [currentItem, setCurrentItem] = useState<Item|undefined>(undefined);
    const [updateCount, setUpdateCount] = useState(0);

    const [confirmAddOpen, setConfirmAddOpen] = React.useState(false);
    const [confirmRemoveOpen, setConfirmRemoveOpen] = React.useState(false);
    const [confirmItemId, setConfirmItemId] = React.useState("");
    const [confirmLocation, setConfirmLocation] = React.useState("");
    const [confirmStatus, setConfirmStatus] = React.useState("");
    const [confirmQuantity, setConfirmQuantity] = React.useState(0);
    
    const handleClickAddIcon = () => {
      items.length > 0 ? setConfirmItemId(items[0].id) : setConfirmItemId("")
      setConfirmAddOpen(true);
      setConfirmLocation("");
      setConfirmStatus("");
      setConfirmQuantity(0);
    };

    const handleClickRemoveIcon = () => {
      items.length > 0 ? setConfirmItemId(items[0].id) : setConfirmItemId("")
      setConfirmRemoveOpen(true);
      setConfirmLocation("");
      setConfirmStatus("");
      setConfirmQuantity(0);
    };
  
    const handleChange = (event: SelectChangeEvent) => {
      setConfirmItemId(event.target.value);
      console.log(event.target.value);
    };

    const handleConfirmAddClose = () => {
      let item = items.find(i => i.id === confirmItemId)
      if (item && confirmLocation && confirmStatus && confirmQuantity) {
        
        fetch(
          'http://localhost:8080/api/addInventory',
          {
            method: 'POST', 
            headers: {  
              'Authorization': 'Bearer 1234'  
            },
            body: JSON.stringify({
              id: "",
              item: {
                id: item.id,
	              sku: item.sku,
	              description: item.description,
              },
              location: confirmLocation,
              status: confirmStatus,
              quantity: confirmQuantity
            })
          },
          )
          .then((response) => {
            console.log("Refresh" + response.status);
            setUpdateCount(updateCount+1)
          })
          .catch((err) => {
              console.log(err.message);
          });
      }
      setConfirmAddOpen(false);
    };
    const handleConfirmRemoveClose = () => {
      setConfirmRemoveOpen(false);
    };

    const handleLocationChange = (event: React.ChangeEvent<HTMLInputElement>) => {
      setConfirmLocation(event.target.value);
    };

    const handleStatusChange = (event: React.ChangeEvent<HTMLInputElement>) => {
      setConfirmStatus(event.target.value);
    };

    const handleQuantityChange = (event: React.ChangeEvent<HTMLInputElement>) => {
      var i = parseInt(event.target.value);
      setConfirmQuantity(i);
    };

    useEffect(() => {
      fetch(
        'http://localhost:8080/api/items',
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
            setItems(data);
            if (data.length > 0) {
              setCurrentItem(data[0])
            } 
         })
         .catch((err) => {
            console.log(err.message);
         });
   }, []);
  
    return (
      <ThemeProvider theme={mdTheme}>
        <Box sx={{ display: 'flex' }}>
          <CssBaseline />
          <AppBar position="absolute" open={drawerOpen}>
            <Toolbar
              sx={{
                pr: '24px', // keep right padding when drawer closed
              }}
            >
              <IconButton
                edge="start"
                color="inherit"
                aria-label="open drawer"
                onClick={toggleDrawer}
                sx={{
                  marginRight: '36px',
                  ...(drawerOpen && { display: 'none' }),
                }}
              >
                <MenuIcon />
              </IconButton>
              <Typography
                component="h1"
                variant="h6"
                color="inherit"
                noWrap
                sx={{ flexGrow: 1 }}
              >
                Available Inventory
              </Typography>
              <IconButton color="inherit" onClick={handleClickAddIcon}>
                <CreateIcon />
              </IconButton>
              <IconButton color="inherit" onClick={handleClickRemoveIcon}>
                <DeleteIcon />
              </IconButton>

              <Dialog open={confirmAddOpen} onClose={handleConfirmAddClose}>
                <DialogTitle>Add Inventory</DialogTitle>
                <DialogContent>
                  <InputLabel id="demo-simple-select-label">SKU</InputLabel>
                  <Select
                    labelId="demo-simple-select-label"
                    id="demo-simple-select"
                    value={confirmItemId}
                    label="Age"
                    onChange={handleChange}
                  >
                    {items.map((item) => {
                      return (
                        <MenuItem key={item.id} value={item.id}>{item.description}</MenuItem>
                      )
                    })}
                  </Select>
                  <TextField
                    autoFocus
                    margin="dense"
                    id="addLocation"
                    label="Location"
                    type="string"
                    fullWidth
                    variant="standard"
                    onChange={handleLocationChange}
                  />
                  <TextField
                    autoFocus
                    margin="dense"
                    id="addStatus"
                    label="Status"
                    type="string"
                    fullWidth
                    variant="standard"
                    onChange={handleStatusChange}
                  />
                  <TextField
                    autoFocus
                    margin="dense"
                    id="addQuantity"
                    label="Quantity"
                    type="Number"
                    fullWidth
                    variant="standard"
                    onChange={handleQuantityChange}
                  />
                </DialogContent>
                <DialogActions>
                  <Button onClick={handleConfirmAddClose}>Cancel</Button>
                  <Button onClick={handleConfirmAddClose}>Add</Button>
                </DialogActions>
              </Dialog>


              <Dialog open={confirmRemoveOpen} onClose={handleConfirmRemoveClose}>
                <DialogTitle>Remove Inventory</DialogTitle>
                <DialogContent>
                  <InputLabel id="demo-simple-select-label">SKU</InputLabel>
                  <Select
                    labelId="demo-simple-select-label"
                    id="demo-simple-select"
                    value={confirmItemId}
                    label="Age"
                    onChange={handleChange}
                  >
                    {items.map((item) => {
                      return (
                        <MenuItem key={item.id} value={item.id}>{item.description}</MenuItem>
                      )
                    })}
                  </Select>
                  <TextField
                    autoFocus
                    margin="dense"
                    id="removeLocation"
                    label="Location"
                    type="string"
                    fullWidth
                    variant="standard"
                  />
                  <TextField
                    autoFocus
                    margin="dense"
                    id="removeStatus"
                    label="Status"
                    type="string"
                    fullWidth
                    variant="standard"
                  />
                  <TextField
                    autoFocus
                    margin="dense"
                    id="removeQuantity"
                    label="Quantity"
                    type="Number"
                    fullWidth
                    variant="standard"
                    onChange={handleQuantityChange}
                  />
                </DialogContent>
                <DialogActions>
                  <Button onClick={handleConfirmRemoveClose}>Cancel</Button>
                  <Button onClick={handleConfirmRemoveClose}>Remove</Button>
                </DialogActions>
              </Dialog>

            </Toolbar>
          </AppBar>
          <Drawer variant="permanent" open={drawerOpen}>
            <Toolbar
              sx={{
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'flex-end',
                px: [1],
              }}
            >
              <IconButton onClick={toggleDrawer}>
                <ChevronLeftIcon />
              </IconButton>
            </Toolbar>
            <Divider />
            <List component="nav">
              <Items
                items={items}
                selected={currentItem}
                setSelected={setCurrentItem}
               />
            </List>
          </Drawer>
          <Box
            component="main"
            sx={{
              backgroundColor: (theme) =>
                theme.palette.mode === 'light'
                  ? theme.palette.grey[100]
                  : theme.palette.grey[900],
              flexGrow: 1,
              height: '100vh',
              overflow: 'auto',
            }}
          >
            <Toolbar />
            <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>  
              <Grid container spacing={3}>
                {/* Item Details */}
                <Grid item xs={12} md={4} lg={3}>
                  <Paper
                    sx={{
                      p: 2,
                      display: 'flex',
                      flexDirection: 'column',
                      height: 240,
                    }}
                  >
                    <ItemDetails current={currentItem} items={items}/>
                  </Paper>
                </Grid>                
                {/* Available Inventory */}
                <Grid item xs={12} md={8} lg={9}>
                  <Paper sx={{ p: 2, display: 'flex', flexDirection: 'column' }}>
                    <InventoryDetails current={currentItem} key={"InventoryDetails-"+updateCount}/>
                  </Paper>
                </Grid>
              </Grid>
            </Container>
          </Box>
        </Box>
      </ThemeProvider>
    );
  }

export default function Dashboard() {
    return <DashboardContent />;
  }