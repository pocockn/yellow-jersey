import * as React from 'react';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Paper from '@mui/material/Paper';
import {Button} from "@mui/material";
import EditIcon from '@mui/icons-material/Edit';


export default function EventsTable({events, handleRoute}) {
    return (
        <TableContainer component={Paper}>
            <Table sx={{minWidth: 650}} size="small" aria-label="a dense table">
                <TableHead>
                    <TableRow>
                        <TableCell>Name</TableCell>
                        <TableCell align="right">Participants</TableCell>
                        <TableCell align="right">Created At</TableCell>
                        <TableCell align="right">Edit</TableCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                    {Array.isArray(events)
                        ? events.map((row) => (
                        <TableRow
                            key={row.name}
                            sx={{'&:last-child td, &:last-child th': {border: 0}}}
                        >
                            <TableCell component="th" scope="row">
                                {row.name}
                            </TableCell>
                            <TableCell align="right">{Array.isArray(row.users) ? row.users.length : ""}</TableCell>
                            <TableCell align="right">{row.created_at}</TableCell>
                            <TableCell align="right">
                                <Button size="small" variant="outlined"
                                        onClick={() => handleRoute(row.id)} endIcon={<EditIcon/>}>
                                </Button>
                            </TableCell>
                        </TableRow>
                    )) : null }
                </TableBody>
            </Table>
        </TableContainer>
    );
}