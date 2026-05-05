import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Plus } from 'lucide-react';
import type { Document as DocType, LineItem as LineItemType } from '../types';

interface LineItemsTableProps {
  doc: DocType;
  lineItems: LineItemType[];
  isEditable: boolean;
  onUpdateLineItem: (index: number, key: keyof LineItemType, value: unknown) => void;
  onAddLineItem: () => void;
}

export function LineItemsTable({
  doc,
  lineItems,
  isEditable,
  onUpdateLineItem,
  onAddLineItem,
}: LineItemsTableProps) {
  return (
    <Card>
      <CardHeader>
        <div className="flex justify-between items-center">
          <CardTitle className="text-lg">Line Items</CardTitle>
          {isEditable && (
            <Button onClick={onAddLineItem} size="sm" variant="outline">
              <Plus className="mr-2 h-4 w-4" />
              Add Item
            </Button>
          )}
        </div>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Description</TableHead>
              <TableHead>Quantity</TableHead>
              <TableHead>Unit Price</TableHead>
              <TableHead>Total</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {lineItems?.length > 0 ? (
              lineItems?.map((item, index) => (
                <TableRow key={index}>
                  <TableCell>
                    {isEditable ? (
                      <Input
                        value={item.description}
                        onChange={(e) => onUpdateLineItem(index, 'description', e.target.value)}
                        className="text-sm"
                      />
                    ) : (
                      item.description
                    )}
                  </TableCell>
                  <TableCell>
                    {isEditable ? (
                      <Input
                        type="number"
                        value={item.quantity}
                        onChange={(e) => onUpdateLineItem(index, 'quantity', parseFloat(e.target.value) || 0)}
                        className="text-sm"
                      />
                    ) : (
                      item.quantity
                    )}
                  </TableCell>
                  <TableCell>
                    {isEditable ? (
                      <Input
                        type="number"
                        value={item.price}
                        onChange={(e) => onUpdateLineItem(index, 'price', parseFloat(e.target.value) || 0)}
                        className="text-sm"
                      />
                    ) : (
                      `${doc.currency} ${item.price}`
                    )}
                  </TableCell>
                  <TableCell>
                    {`${doc.currency} ${item.total}`}
                  </TableCell>
                </TableRow>
              ))
            ) : (
              <TableRow>
                <TableCell colSpan={4} className="text-center text-gray-500">
                  No line items available.
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  );
}
